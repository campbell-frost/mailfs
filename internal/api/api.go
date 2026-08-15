package api

import (
	"net/http"
	"time"

	"github.com/campbell-frost/mailfs/internal/vault"
)

type Server struct {
	vault   *vault.Vault
	mux     *http.ServeMux
	handler http.Handler
}

func New(v *vault.Vault) *Server {
	s := &Server{
		vault: v,
		mux:   http.NewServeMux(),
	}

	s.routes()
	s.handler = cors(s.mux)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /upload/", s.UploadHandler)
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

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
