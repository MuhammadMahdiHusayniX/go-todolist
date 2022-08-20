package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MuhammadMahdiHusayniX/go-todolist/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// var conf *oauth2.Config

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func init() {
	models.Setup()
}

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

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
	// session := sessions.Default(c)
	// session.Set("state", "state")
	// err := session.Save()
	// if err != nil {
	// 	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
	// 	return
	// }
	// link := getLoginURL(state)
	// c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link})
}

func AuthHandler(c *gin.Context) {
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
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	c.HTML(http.StatusOK, "authenticated.tmpl", gin.H{
		"user": user,
	})
}

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
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
	r.Use(sessions.Sessions("goquestsession", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/", IndexHandler)
	r.GET("/login", LoginHandler)
	r.GET("/auth/google/callback", AuthHandler)

	authorized := r.Group("/todo")
	authorized.Use(AuthorizeRequest())
	{
		authorized.GET("", RetrieveAll)
		authorized.GET("/complete/:id", MarkComplete)
		authorized.POST("", AddTodo)
		authorized.DELETE("/:id", DeleteTodo)
	}
	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}
