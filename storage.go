package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TodoStorage interface {
}

type StorageImpl struct {
	connection *sql.DB
}

func NewStorage(dbType, dbName string) (StorageImpl, error) {
	connection, err := sql.Open(dbType, dbName)
	if err != nil {
		return StorageImpl{}, err
	}
	log.Printf("Connected to database %v", connection)
	return StorageImpl{
		connection: connection,
	}, nil
}
