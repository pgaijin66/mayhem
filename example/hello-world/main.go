package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "world"})
}

func main() {
	router := gin.New()
	router.GET("/hello", hello)

	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})
	})

	router.POST("/users", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	})

	router.Run(":9090")
}
