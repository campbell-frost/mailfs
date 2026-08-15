package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/campbell-frost/mailfs/internal/handler"
)

const addr = "localhost:1738"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload/", handler.UploadHandler)
	fmt.Printf("server running on http://%v\n", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}
