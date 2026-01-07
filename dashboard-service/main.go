package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ticketServiceURL := os.Getenv("TICKET_SERVICE_URL")
	waitingRoomURL := os.Getenv("WAITING_ROOM_URL")

	// shared client + timeouts
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	r := gin.Default()

	r.GET("/dashboard", func(c *gin.Context) {
		// ---- ticket-service stats (only sold tickets now) ----
		ticketResp, err := httpClient.Get(ticketServiceURL + "/stats")
		if err != nil {
			c.JSON(500, gin.H{"error": "ticket service unavailable"})
			return
		}
		defer ticketResp.Body.Close()

		var ticketStats struct {
			SoldTickets int `json:"sold_tickets"`
		}
		if err := json.NewDecoder(ticketResp.Body).Decode(&ticketStats); err != nil {
			c.JSON(500, gin.H{"error": "failed to decode ticket stats"})
			return
		}

		// ---- waiting-room stats (active + waiting) ----
		waitResp, err := httpClient.Get(waitingRoomURL + "/stats")
		if err != nil {
			c.JSON(500, gin.H{"error": "waiting room unavailable"})
			return
		}
		defer waitResp.Body.Close()

		var waitingStats struct {
			WaitingUsers int `json:"waiting_users"`
			ActiveBuyers int `json:"active_buyers"`
		}
		if err := json.NewDecoder(waitResp.Body).Decode(&waitingStats); err != nil {
			c.JSON(500, gin.H{"error": "failed to decode waiting stats"})
			return
		}

		// ---- merged response ----
		c.JSON(200, gin.H{
			"sold_tickets":  ticketStats.SoldTickets,
			"active_buyers": waitingStats.ActiveBuyers,
			"waiting_users": waitingStats.WaitingUsers,
		})
	})

	r.Run(":8083")
}
