package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	waitingRoomURL := os.Getenv("WAITING_ROOM_URL")
	dashboardURL := os.Getenv("DASHBOARD_SERVICE_URL")

	// ✅ Reusable HTTP client with timeouts + pooling
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// route to buy tickets
	r.GET("/buy", func(c *gin.Context) {
		resp, err := httpClient.Get(waitingRoomURL + "/enter")
		if err != nil {
			c.JSON(500, gin.H{"error": "waiting room down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// route to dashboard service
	r.GET("/dashboard", func(c *gin.Context) {
		resp, err := httpClient.Get(dashboardURL + "/dashboard")
		if err != nil {
			c.JSON(500, gin.H{"error": "dashboard service down"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.Run(":8080")
}
