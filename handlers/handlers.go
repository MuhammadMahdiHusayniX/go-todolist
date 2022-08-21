package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
