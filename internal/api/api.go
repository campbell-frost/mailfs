package api

import (
	"net/http"
	"time"

	"github.com/campbell-frost/mailfs/internal/vault"
)

type Server struct {
	vault *vault.Vault
	mux   *http.ServeMux
}

func New(v *vault.Vault) *Server {
	s := &Server{
		vault: v,
		mux:   http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/upload/", s.UploadHandler)
	s.mux.HandleFunc("GET /api/list/", s.ListHandler)
	s.mux.HandleFunc("GET /api/download/{id}/", s.DownloadHandler)
	s.mux.HandleFunc("DELETE /api/delete/{id}/", s.DeleteHandler)
	s.mux.HandleFunc("GET /api/status/{id}/", s.StatusHandler)
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
