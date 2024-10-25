package main

import (
	"encoding/json"
	"html/template"
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

func (s ServerImpl) newTodoItemPage(writer http.ResponseWriter, r *http.Request) {
	http.ServeFile(writer, r, "new_todo_item.html")
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

func (s ServerImpl) listTodoPageHandler(writer http.ResponseWriter, r *http.Request) {
	// construct template on the fly - allow us to change the template
	// while the service is running
	const templateFilename = "todo_list.html"
	log.Printf("Constructing template from file %s", templateFilename)
	// new template
	tmpl, err := template.ParseFiles(templateFilename)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Template can't be constructed: %v", err)
		return
	}

	todoList, err := s.todoStorage.ReadTodoList()
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		_, err := writer.Write([]byte("Unable to retrieve list of TODOs"))
		if err != nil {
			log.Printf("Unable to retrieve list of TODOs: %v", err)
		}
		return
	}
	log.Printf("Application template for %d data records", len(todoList))

	// apply template
	err = tmpl.Execute(writer, todoList)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func (s ServerImpl) addNewTodoItem(writer http.ResponseWriter, r *http.Request) {
	dueTo := r.FormValue("due_to")
	finished := r.FormValue("finished")
	priority := r.FormValue("priority")
	subject := r.FormValue("subject")
	details := r.FormValue("details")
	log.Println("Adding new TODO item", dueTo, finished, priority, subject, details)
	s.todoStorage.AddNewTodoItem(dueTo, finished, priority, subject, details)
	http.ServeFile(writer, r, "index.html")
}

// startServer starts HTTP server that provides all static and dynamic data
func (s ServerImpl) Serve(port uint) {
	log.Printf("Starting server on port %d", port)
	// HTTP pages
	http.HandleFunc("/", s.indexPageHandler)
	http.HandleFunc("/list-todo", s.listTodoPageHandler)
	http.HandleFunc("/new-todo-item", s.newTodoItemPage)

	// REST API endpoints
	http.HandleFunc("/todo", s.returnListOfTodosAsJSON)
	http.HandleFunc("/add-new-todo-item", s.addNewTodoItem)

	// start the server
	// TODO: use proper port number!!!
	http.ListenAndServe(":8080", nil)
}
