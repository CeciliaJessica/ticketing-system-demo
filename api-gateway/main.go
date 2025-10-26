package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	waitingRoomURL := os.Getenv("WAITING_ROOM_URL")

	r := gin.Default()

	r.GET("/buy", func(c *gin.Context) {
		resp, err := http.Get(waitingRoomURL + "/enter")
		if err != nil {
			c.JSON(500, gin.H{"error": "waiting room down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
	})

	r.Run(":8080")
}
