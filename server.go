package main

import (
	"encoding/json"
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

func (s ServerImpl) returnListOfTodosAsJSON(writer http.ResponseWriter, r *http.Request) {
	todoList, err := s.todoStorage.ReadTodoList()
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Storage error: %v", err)
		_, err := writer.Write([]byte("Unable to retrieve list TODOs"))
		if err != nil {
			log.Printf("Unable to retrieve list of TODOs: %v", err)
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todoList)
}

// startServer starts HTTP server that provides all static and dynamic data
func (s ServerImpl) Serve(port uint) {
	log.Printf("Starting server on port %d", port)
	// HTTP pages
	http.HandleFunc("/", s.indexPageHandler)

	// REST API endpoints
	http.HandleFunc("/todo", s.returnListOfTodosAsJSON)

	// start the server
	// TODO: use proper port number!!!
	http.ListenAndServe(":8080", nil)
}
