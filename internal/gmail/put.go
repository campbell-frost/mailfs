package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Chunk struct {
	FileID string
	Seq    int
	Sha256 string
	Size   int
	Data   []byte
}

type ref struct {
	Mailbox     string `json:"mailbox"`
	UID         uint32 `json:"uid"`
	UIDValidity uint32 `json:"uid_validity"`
	MessageID   string `json:"message_id"`
}

// Put takes a chunk and appends it to the mailbox.  We then create an opaque
// reference object and return that to the index for future retrieval
func (c *Client) Put(ctx context.Context, ch Chunk) ([]byte, error) {
	mID := messageID(ch.FileID, ch.Seq)
	raw := encode(c.cfg.Username, mID, ch, time.Now())
	var r ref

	if err := c.do(ctx, func(cn *conn) error {
		cmd := cn.cl.Append(c.cfg.Mailbox, int64(len(raw)), nil)

		if _, err := cmd.Write(raw); err != nil {
			return err
		}

		if err := cmd.Close(); err != nil {
			return err
		}

		data, err := cmd.Wait()
		if err != nil {
			return err
		}

		if data.UID == 0 {
			return fmt.Errorf("append %s: no APPENDUID in response", mID)
		}

		r = ref{
			Mailbox:     c.cfg.Mailbox,
			UID:         uint32(data.UID),
			UIDValidity: data.UIDValidity,
			MessageID:   mID,
		}

		if r.UIDValidity == 0 {
			r.UIDValidity = cn.uidValidity
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return json.Marshal(r)
}

const (
	// MIME caps a base64 line at 76 chars, and base64 turns 3 bytes into 4,
	// so 57 is the largest input slice that fits
	charsPerLine = 76
	rawPerLine   = charsPerLine / 4 * 3

	headerSize = 512
	crlfSize   = 2
)

const (
	hdrFileID = "X-MAILFS-File-Id"
	hdrSeq    = "X-MAILFS-Seq"
	hdrSize   = "X-MAILFS-Size"
	hdrSha256 = "X-MAILFS-Sha256"
)

func encode(username string, mID string, c Chunk, now time.Time) []byte {
	var buf bytes.Buffer

	// one line per rawPerLine bytes of input + header
	buf.Grow(len(c.Data)/rawPerLine*(charsPerLine+crlfSize) + headerSize)

	hdr := func(k, v string) { fmt.Fprintf(&buf, "%s: %s\r\n", k, v) }

	hdr("From", username)
	hdr("To", username)
	hdr("Subject", fmt.Sprintf("mailfs/%s/%d", c.FileID, c.Seq))
	hdr("Message-ID", mID)
	hdr("Date", now.Format(time.RFC1123Z))
	hdr("MIME-Version", "1.0")
	hdr(hdrFileID, c.FileID)
	hdr(hdrSeq, strconv.Itoa(c.Seq))
	hdr(hdrSize, strconv.Itoa(c.Size))
	hdr(hdrSha256, c.Sha256)
	hdr("Content-Type", "application/octet-stream")
	hdr("Content-Transfer-Encoding", "base64")
	buf.WriteString("\r\n")

	line := make([]byte, base64.StdEncoding.EncodedLen(rawPerLine))
	for i := 0; i < len(c.Data); i += rawPerLine {
		end := min(i+rawPerLine, len(c.Data))
		base64.StdEncoding.Encode(line, c.Data[i:end])
		buf.Write(line[:base64.StdEncoding.EncodedLen(end-i)])
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}

func messageID(fileID string, seq int) string {
	return fmt.Sprintf("<mailfs.%s.%d@mailfs.local>", fileID, seq)
}
