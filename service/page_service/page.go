package page_service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorPage(c *gin.Context, errorMessage string) {
	c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
}
