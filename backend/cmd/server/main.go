package main

import (
	"log"
	"vibeproxyplus/backend/internal/api"
)

func main() {
	srv, err := api.NewServer()
	if err != nil {
		log.Fatalf("init server: %v", err)
	}

	log.Printf("backend listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
