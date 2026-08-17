package main

import (
	"log"
	"os"

	"github.com/campbell-frost/mailfs/internal/api"
	"github.com/campbell-frost/mailfs/internal/gmail"
	"github.com/campbell-frost/mailfs/internal/index"
	"github.com/campbell-frost/mailfs/internal/vault"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load env: %v", err)
	}

	idx, err := index.New("mailfs.db")
	if err != nil {
		log.Fatalf("failed to create index: %v", err)
	}
	defer idx.Close()

	c, err := gmail.NewClient(gmail.Config{
		Username: os.Getenv("MAILFS_IMAP_USER"),
		Password: os.Getenv("MAILFS_IMAP_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("gmail: failed to create client: %v", err)
	}
	defer c.Close()

	v := vault.New(idx, c, "tmp")

	s := api.New(v)

	addr := "localhost:1738"
	log.Printf("server running on http://%s", addr)
	log.Fatal(s.Start(addr))
}
