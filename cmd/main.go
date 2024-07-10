package main

import (
	"github.com/gin-gonic/gin"
	"log"
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
	err := r.Run()
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
