package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	//  Simple demo endpoint
	r.GET("/demo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// Start the server
	r.Run()
}
