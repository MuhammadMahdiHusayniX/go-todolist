package main

import (
	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id        int
	Text      string
	CreatedBy string
	Completed bool
}

var Todos = []Todo{
	{
		Id:        1,
		Text:      "Optimize code",
		CreatedBy: "Mahdi",
		Completed: false,
	},
	{
		Id:        2,
		Text:      "Queueing system for api call",
		CreatedBy: "Mahdi",
		Completed: false,
	},
}

func AddTodo(c *gin.Context) {
	var req Todo
	c.BindJSON(&req)
	Todos = append(Todos, req)
	c.JSON(200, Todos)
}

// more research on how to get id
func DeleteTodo(c *gin.Context) {
}

func RetrieveAll(c *gin.Context) {
	c.JSON(200, Todos)
}

func main() {
	r := gin.Default()
	r.POST("/todo", AddTodo)
	r.DELETE("/todo/:Id", DeleteTodo)
	r.GET("/todo", RetrieveAll)
	r.Run() // listen and serve on 0.0.0.0:8080
}
