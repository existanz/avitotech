package main

import (
	"avitotech/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("http server error: %s", err)
	}
}
