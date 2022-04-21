package main

import (
	log "github.com/sirupsen/logrus"
	server "kilimanjaro-api/server"
)

func main() {

	s, err := server.NewServer("/api")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on localhost")

	err = s.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
