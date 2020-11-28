package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":8000",
		Handler: http.FileServer(http.Dir(cwd)),
	}
	log.Fatal(srv.ListenAndServe())
}
