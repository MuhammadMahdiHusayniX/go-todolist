package todo_service

import (
	"net/http"
	"strconv"

	"github.com/MuhammadMahdiHusayniX/go-todolist/models"
	"github.com/gin-gonic/gin"
)

func AddTodo(c *gin.Context) {
	todo := models.Todo{}

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	todoId, err := models.AddTodo(todo.Task, todo.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"todoId": todoId,
	})
}

func DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	deleteErr := models.DeleteTodo(id)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, deleteErr)
	}

	c.JSON(http.StatusOK, gin.H{
		"Successfully delete todo id: ": id,
	})
}

func RetrieveAll(c *gin.Context) {
	todos, err := models.GetTodos()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, todos)
}

func MarkComplete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	markCompleteError := models.MarkComplete(id)
	if markCompleteError != nil {
		c.JSON(http.StatusBadRequest, markCompleteError)
	}

	c.JSON(http.StatusOK, gin.H{
		"Successfully mark complete for id: ": id,
	})
}
