package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type TodoStorage interface {
	ReadTodoList() ([]TODO, error)
	AddNewTodoItem(dueTo, finished, priority, subject, details string) error
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

func (s StorageImpl) ReadTodoList() ([]TODO, error) {
	rows, err := s.connection.Query("SELECT id, due_to, finished, priority, subject, details FROM todo ORDER BY id")
	if err != nil {
		return []TODO{}, err
	}
	defer rows.Close()

	var todoList []TODO

	for rows.Next() {
		var id int
		var due_to string
		var finished bool
		var priority int
		var subject string
		var details string

		err := rows.Scan(&id, &due_to, &finished, &priority, &subject, &details)
		if err != nil {
			return todoList, err
		}
		todoList = append(todoList, TODO{
			ID:       id,
			DueTo:    due_to,
			Finished: finished,
			Priority: priority,
			Subject:  subject,
			Details:  details,
		})
	}

	return todoList, nil
}

func (s StorageImpl) AddNewTodoItem(dueTo, finished, priority, subject, details string) error {
	finishedB := finished == "Yes" || finished == "y"
	priorityI, err := strconv.Atoi(priority)
	if err != nil {
		return err
	}
	statement, err := s.connection.Prepare("INSERT INTO todo(due_to, finished, priority, subject, details) VALUES (?, ?, ?, ?, ?)")
	_, err = statement.Exec(dueTo, finishedB, priorityI, subject, details)
	return err
}
