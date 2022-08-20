package main

import (
	"strconv"

	"github.com/MuhammadMahdiHusayniX/go-todolist/models"
	"github.com/gin-gonic/gin"
)

func init() {
	models.Setup()
}

func AddTodo(c *gin.Context) {
	todo := models.Todo{}

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(400, err)
	}

	todoId, err := models.AddTodo(todo.Task, todo.CreatedBy)
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(200, gin.H{
		"todoId": todoId,
	})
}

func DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, err)
	}

	deleteErr := models.DeleteTodo(id)
	if deleteErr != nil {
		c.JSON(400, deleteErr)
	}

	c.JSON(200, gin.H{
		"Successfully delete todo id: ": id,
	})
}

func RetrieveAll(c *gin.Context) {
	todos, err := models.GetTodos()
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(200, todos)
}

func MarkComplete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, err)
	}

	markCompleteError := models.MarkComplete(id)
	if markCompleteError != nil {
		c.JSON(400, markCompleteError)
	}

	c.JSON(200, gin.H{
		"Successfully mark complete for id: ": id,
	})
}

func main() {
	r := gin.Default()
	r.GET("/todo", RetrieveAll)
	r.GET("/todo/complete/:id", MarkComplete)
	r.POST("/todo", AddTodo)
	r.DELETE("/todo/:id", DeleteTodo)
	r.Run() // listen and serve on 0.0.0.0:8080
}
