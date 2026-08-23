package gmail

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strconv"

	"github.com/emersion/go-imap/v2"
)

func (c *Client) Get(ctx context.Context, chunkRef []byte) ([]byte, error) {
	r, err := parseRef(chunkRef)
	if err != nil {
		return nil, err
	}

	var bodyRaw []byte
	if err := c.do(ctx, func(cn *conn) error {
		uids, err := c.locate(cn, r)
		fmt.Println(uids, "uids from fetch")
		if err != nil {
			return err
		}

		cmd := cn.cl.Fetch(uids, &imap.FetchOptions{
			BodySection: []*imap.FetchItemBodySection{{Peek: true}},
		})

		msgs, err := cmd.Collect()
		if err != nil {
			return err
		}
		fmt.Println(msgs, "msgs from fetch")

		if len(msgs) == 0 || len(msgs[0].BodySection) == 0 {
			return fmt.Errorf("no body section")
		}
		bodyRaw = msgs[0].BodySection[0].Bytes

		return nil
	}); err != nil {
		return nil, err
	}

	data, err := parseMessages(bodyRaw)
	if err != nil {
		return nil, err
	}

	return data, nil
}

var errNotFound = errors.New("chunk not found")

func (c *Client) locate(cn *conn, r ref) (imap.UIDSet, error) {
	if r.UID != 0 && r.UIDValidity == cn.uidValidity{
		return imap.UIDSetNum(imap.UID(r.UID)), nil
	}
	if r.MessageID == "" {
		return nil, errNotFound
	}
	sc := &imap.SearchCriteria{
		Header: []imap.SearchCriteriaHeaderField{
			{ Key: "Message-Id", Value: r.MessageID},
		},
	}
	data, err := cn.cl.UIDSearch(sc, &imap.SearchOptions{ReturnAll: true}).Wait()
	if err != nil {
		return nil, err
	}

	uids := data.AllUIDs()
	if len(uids) == 0 {
		return nil, errNotFound
	}

	return imap.UIDSetNum(uids[0]), nil
}


func parseMessages(raw []byte) ([]byte, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	fileID := msg.Header.Get(hdrFileID)
	if fileID == "" {
		return nil, fmt.Errorf("no file id")
	}

	_, err = strconv.Atoi(msg.Header.Get(hdrSeq))
	if err != nil {
		return nil, fmt.Errorf("bad %s: %w", hdrSeq, err)
	}
	data, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, msg.Body))
	if err != nil {
		return nil, fmt.Errorf("bad body: %w", err)
	}

	want := msg.Header.Get(hdrSha256)
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if want != "" && want != got {
		return nil, fmt.Errorf("chechsum mismatch: want %s got %s", want, got)
	}

	return data, nil
}
