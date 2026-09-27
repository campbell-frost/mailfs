package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/campbell-frost/mailfs/internal/index"
	"golang.org/x/oauth2"
)

const (
	sessionCookie = "session"
	oauthCookie   = "oauth"
	sessionTTL    = 30 * 24 * time.Hour
)

type userResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *Server) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := rand.Text()
	verifier := oauth2.GenerateVerifier()

	http.SetCookie(w, &http.Cookie{
		Name:     oauthCookie,
		Value:    state + "." + verifier,
		Path:     "/auth/google",
		MaxAge:   int((10 * time.Minute).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, s.oauth.AuthURL(state, verifier), http.StatusFound)
}

func (s *Server) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if e := r.URL.Query().Get("error"); e != "" {
		log.Printf("google callback error: %s", e)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	c, err := r.Cookie(oauthCookie)
	if err != nil {
		http.Error(w, "missing oauth cookie", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   oauthCookie,
		Path:   "/auth/google",
		MaxAge: -1,
	})

	state, verifier, _ := strings.Cut(c.Value, ".")

	// constant time so the comparison doesn't leak how much of the state matched
	if state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(r.URL.Query().Get("state"))) != 1 {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	// get user info from the code
	info, err := s.oauth.Exchange(
		r.Context(),
		r.URL.Query().Get("code"),
		verifier,
	)
	if err != nil {
		log.Printf("google callback failed: %v", err)
		http.Error(w, "failed to sign in with google", http.StatusBadGateway)
		return
	}

	u, err := s.idx.User.Upsert(r.Context(), index.User{
		Sub:   info.Sub,
		Email: info.Email,
		Name:  info.Name,
	})
	if err != nil {
		log.Printf("failed to upsert user %s: %v", info.Email, err)
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

	// create a new session
	token := rand.Text()
	now := time.Now()
	err = s.idx.Session.Create(r.Context(), index.Session{
		TokenHash: hashToken(token),
		UserID:    u.ID,
		CreatedAt: now,
		ExpiresAt: now.Add(sessionTTL),
	})
	if err != nil {
		log.Printf("failed to create session for %s: %v", u.Email, err)
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	log.Printf("user %d (%s) signed in", u.ID, u.Email)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) SessionHandler(w http.ResponseWriter, r *http.Request) {
	u, err := s.currentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{ID: u.ID, Email: u.Email, Name: u.Name})
}

// deleting the session effectivly logs the user out
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := s.idx.Session.Delete(r.Context(), hashToken(c.Value)); err != nil {
			log.Printf("failed to delete session: %v", err)
		}
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Path: "/", MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) requireAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := s.currentUser(r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

// currentUser resolves the session cookie to a user. If we fail, the user is not authenticated
func (s *Server) currentUser(r *http.Request) (index.User, error) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return index.User{}, err
	}

	sess, err := s.idx.Session.Get(r.Context(), hashToken(c.Value))
	if err != nil {
		return index.User{}, err
	}
	if time.Now().After(sess.ExpiresAt) {
		s.idx.Session.Delete(context.WithoutCancel(r.Context()), sess.TokenHash)
		return index.User{}, index.ErrNotFound
	}

	return s.idx.User.Get(r.Context(), sess.UserID)
}

func hashToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}
