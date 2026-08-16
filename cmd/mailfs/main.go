package main

import (
	"log"

	"github.com/campbell-frost/mailfs/internal/api"
	"github.com/campbell-frost/mailfs/internal/gmail"
	"github.com/campbell-frost/mailfs/internal/index"
	"github.com/campbell-frost/mailfs/internal/vault"
)

const GMAIL_CHUNK_SIZE = 17 << 20

func main() {
	i, err := index.New("mailfs.db")
	if err != nil {
		log.Fatalf("failed to create index: %e", err)
	}
	defer i.Close()

	g := gmail.New(GMAIL_CHUNK_SIZE)

	v := vault.New(i, g, "tmp")

	s := api.New(v)

	addr := "localhost:1738"
	log.Printf("server running on http://%v\n", addr)
	log.Fatal(s.Start(addr))
}
