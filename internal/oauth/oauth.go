package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

type OAuth struct {
	cfg *oauth2.Config
}

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func New(c Config) *OAuth {
	return &OAuth{cfg: &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}}
}

type UserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (o *OAuth) AuthURL(state, verifier string) string {
	return o.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
}

func (o *OAuth) Exchange(ctx context.Context, code, verifier string) (UserInfo, error) {
	tok, err := o.cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return UserInfo{}, fmt.Errorf("exchange code: %w", err)
	}

	res, err := o.cfg.Client(ctx, tok).Get(userInfoURL)
	if err != nil {
		return UserInfo{}, fmt.Errorf("get userinfo: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf("get userinfo: bad status %d", res.StatusCode)
	}

	var u UserInfo
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return UserInfo{}, fmt.Errorf("decode userinfo: %w", err)
	}
	return u, nil
}
