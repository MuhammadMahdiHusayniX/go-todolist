package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

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

	credsJson, jsonErr := os.Open("conf/creds.json")
	if jsonErr != nil {
		log.Fatal("unable to read credentials: ", jsonErr)
	}

	defer credsJson.Close()

	byteValue, _ := io.ReadAll(credsJson)

	var credentials models.Creds
	unMarshalError := json.Unmarshal(byteValue, &credentials)
	if unMarshalError != nil {
		log.Fatal("Failed to decode json: ", unMarshalError)
	}

	gothic.Store = store
	goth.UseProviders(
		google.New(credentials.Cid, credentials.Csecret, credentials.Ccallback),
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
