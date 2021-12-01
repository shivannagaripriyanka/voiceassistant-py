package main

import (
	"fmt"
	"net/http"

	// "voiceassistant/amodels"
	// "voiceassistant/aroutes"

	"voiceassistant/models"
	"voiceassistant/routes"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.Default()

	router.Use(dbMiddleware(*conn))

	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UsersRegister)
		usersGroup.POST("login", routes.UsersLogin)
	}

	itemsGroup := router.Group("items")
	{
		// itemsGroup.POST("index", authMiddleWare(),routes.ItemsIndex)
		itemsGroup.GET("index", routes.ItemsIndex)
		itemsGroup.POST("create", authMiddleWare(), routes.ItemsCreate)
		itemsGroup.PUT("update", authMiddleWare(), routes.ItemsUpdate)
		itemsGroup.POST("speech", routes.SpeechToText)
		itemsGroup.POST("emotion", routes.Emotion)
		itemsGroup.POST("findloc", routes.FindItem)

	}

	gin.SetMode(gin.ReleaseMode)0
	router.Run(":3000")
}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:1619@localhost:5432/voiceassistant")
	if err != nil || conn == nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]
		fmt.Printf("Bearer (%v) \n", token)
		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
