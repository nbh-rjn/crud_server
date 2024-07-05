package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"encoding/csv"
	"os"
)

type todo struct {
	Id   string `json:"id"`
	Desc string `json:"desc"`
	Done bool   `json:"done"`
}

var todo_list = []todo{
	{Id: "1", Desc: "beego orm", Done: false},
	{Id: "2", Desc: "tableau rest api", Done: false},
	{Id: "3", Desc: "postman ui", Done: false},
}

func create_todo(context *gin.Context) {
	var new_todo todo

	if err := context.BindJSON(&new_todo); err != nil {
		return
	}

	todo_list = append(todo_list, new_todo)

	context.IndentedJSON(http.StatusCreated, new_todo)
}

func read_todos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todo_list)
}

func read_todo(context *gin.Context) {
	id := context.Param("id")
	todo, err := find_todo_by_id(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, todo)
}

func find_todo_by_id(id string) (*todo, error) {
	for i, t := range todo_list {
		if t.Id == id {
			return &todo_list[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func update_todo(context *gin.Context) {
	id := context.Param("id")
	todo, err := find_todo_by_id(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	todo.Done = !todo.Done

	context.IndentedJSON(http.StatusOK, todo)
}

func delete_todo(context *gin.Context) {
	id := context.Param("id")
	for i, t := range todo_list {
		if t.Id == id {
			todo_list = append(todo_list[:i], todo_list[i+1:]...)
			context.IndentedJSON(http.StatusOK, gin.H{"message": "todo deleted"})
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}

func parseCSV(filename string) ([]todo, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var people []todo

	for i, record := range records {
		if i == 0 {
			continue
		}

		var done bool
		if record[2] == "true" {
			done = true
		} else {
			done = false
		}

		todo := todo{
			Id:   record[0],
			Desc: record[1],
			Done: done,
		}

		todo_list = append(todo_list, todo)
	}

	return people, nil
}

func main() {
	router := gin.Default()

	// Create
	router.POST("/todos", create_todo)

	// Read
	router.GET("/todos", read_todos)
	router.GET("/todos/:id", read_todo)

	// Update
	router.PATCH("/todos/:id", update_todo)

	// Delete
	router.DELETE("/todos/:id", delete_todo)

	// Parse CSV data
	todos, err := parseCSV("data.csv")
	if err != nil {
		panic(err)
	}
	todo_list = append(todo_list, todos...)

	router.Run("localhost:9090")
}
