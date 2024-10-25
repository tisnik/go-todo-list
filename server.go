package main

import (
	"log"
	"net/http"
)

type Server interface {
	Serve(port uint)
}

type ServerImpl struct {
	todoStorage TodoStorage
}

func NewServer(todoStorage TodoStorage) Server {
	return ServerImpl{
		todoStorage: todoStorage,
	}
}

func (s ServerImpl) indexPageHandler(writer http.ResponseWriter, r *http.Request) {
	http.ServeFile(writer, r, "index.html")
}

// startServer starts HTTP server that provides all static and dynamic data
func (s ServerImpl) Serve(port uint) {
	log.Printf("Starting server on port %d", port)
	// HTTP pages
	http.HandleFunc("/", s.indexPageHandler)

	// start the server
	// TODO: use proper port number!!!
	http.ListenAndServe(":8080", nil)
}
