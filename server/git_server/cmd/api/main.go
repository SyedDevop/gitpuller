package main

import (
	"git_server/internal/server"

	"github.com/charmbracelet/log"
)

func main() {
	server := server.NewServer()

	log.Infof("Git mock server started @::-> http://localhost%s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("cannot start server: %s", err)
	}
}
