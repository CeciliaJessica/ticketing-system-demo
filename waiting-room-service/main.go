package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ticketServiceUrl := os.Getenv("TICKET_SERVICE_URL")

	r := gin.Default()

	r.GET("/enter", func(c *gin.Context) {
		resp, err := http.Get(ticketServiceUrl + "/active")
		if err != nil {
			c.JSON(500, gin.H{"error": "ticket service unreachable"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		active, ok := data["active_users"].(float64)
		if !ok {
			c.JSON(500, gin.H{"error": "invalid ticket service response"})
			return
		}

		if active < 100 {
			// Let user go to ticket service
			resp2, err := http.Get(ticketServiceUrl + "/tickets")
			if err != nil {
				c.JSON(500, gin.H{"error": "could not reach ticket service"})
				return
			}
			defer resp2.Body.Close()
			c.DataFromReader(resp2.StatusCode, resp2.ContentLength, "application/json", resp2.Body, nil)
		} else {
			c.JSON(429, gin.H{
				"status":  "waiting",
				"message": "Too many users right now, please try again soon.",
			})
		}
	})

	r.Run(":8082")
}
