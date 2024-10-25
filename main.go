package main

import (
	"flag"
	"log"
)

func main() {
	var port uint

	flag.UintVar(&port, "p", 8080, "port for the server (shorthand)")
	flag.UintVar(&port, "port", 8080, "port for the server")

	flag.Parse()

	storage, err := NewStorage("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal("Can not connect to data storage", err)
		return
	}
	log.Println("Connection to storage established", storage)

	log.Println("Starting server on port", port)

}
