package main

import (
	"log"

	"github.com/MuhammadMahdiHusayniX/go-todolist/handlers"
	"github.com/MuhammadMahdiHusayniX/go-todolist/models"
	"github.com/MuhammadMahdiHusayniX/go-todolist/service/auth_service"
	"github.com/MuhammadMahdiHusayniX/go-todolist/service/todo_service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	models.Setup()
}

func main() {
	r := gin.Default()
	token, err := handlers.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := cookie.NewStore([]byte(token))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
	})

	gothic.Store = store

	goth.UseProviders(
		google.New("154301927301-uln6u2mc1ie0ojop9ig0mp0no3bcc776.apps.googleusercontent.com", "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", "http://127.0.0.1:3000/auth/google/callback"),
	)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("gotodosession", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/", handlers.IndexHandler)
	r.GET("/login", auth_service.Login)
	r.GET("/auth/google/callback", auth_service.Auth)

	authorized := r.Group("/todo")
	authorized.Use(auth_service.AuthorizeRequest())
	{
		authorized.GET("", todo_service.RetrieveAll)
		authorized.GET("/complete/:id", todo_service.MarkComplete)
		authorized.POST("", todo_service.AddTodo)
		authorized.DELETE("/:id", todo_service.DeleteTodo)
	}
	r.Run(":3000")
}
