package auth_service

import (
	"log"
	"net/http"

	"github.com/MuhammadMahdiHusayniX/go-todolist/service/page_service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func Auth(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	session.Set("user-id", user.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		page_service.ErrorPage(c, "Error while saving session. Please try again.")
		return
	}

	c.HTML(http.StatusOK, "authenticated.tmpl", gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")
		if v == nil {
			page_service.ErrorPage(c, "Please login.")
			c.Abort()
		}
		c.Next()
	}
}
