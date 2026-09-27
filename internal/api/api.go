package api

import (
	"net/http"
	"time"

	"github.com/campbell-frost/mailfs/internal/index"
	"github.com/campbell-frost/mailfs/internal/oauth"
	"github.com/campbell-frost/mailfs/internal/vault"
)

type Server struct {
	vault *vault.Vault
	idx   *index.Index
	oauth *oauth.OAuth
	mux   *http.ServeMux
}

func New(
	v *vault.Vault,
	idx *index.Index,
	o *oauth.OAuth,
) *Server {
	s := &Server{
		vault: v,
		idx:   idx,
		oauth: o,
		mux:   http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/upload/", s.requireAuth(s.UploadHandler))
	s.mux.HandleFunc("GET /api/list/", s.requireAuth(s.ListHandler))
	s.mux.HandleFunc("GET /api/download/{id}/", s.requireAuth(s.DownloadHandler))
	s.mux.HandleFunc("DELETE /api/delete/{id}/", s.requireAuth(s.DeleteHandler))
	s.mux.HandleFunc("GET /api/status/{id}/", s.requireAuth(s.StatusHandler))

	// auth
	s.mux.HandleFunc("GET /auth/google/login", s.GoogleLoginHandler)
	s.mux.HandleFunc("GET /auth/google/callback", s.GoogleCallbackHandler)
	s.mux.HandleFunc("GET /api/session", s.SessionHandler)
	s.mux.HandleFunc("DELETE /api/session", s.LogoutHandler)
}

func (s *Server) Start(addr string) error {
	svr := &http.Server{
		Addr:              addr,
		Handler:           s,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return svr.ListenAndServe()
}

type fileInfoResponse struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	Chunks    int       `json:"chunks"`
}
