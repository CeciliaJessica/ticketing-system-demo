package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ticketServiceURL := os.Getenv("TICKET_SERVICE_URL")
	waitingRoomURL := os.Getenv("WAITING_ROOM_URL")

	r := gin.Default()

	r.GET("/dashboard", func(c *gin.Context) {
		// fetch ticket-service stats
		ticketResp, err := http.Get(ticketServiceURL + "/stats")
		if err != nil {
			c.JSON(500, gin.H{"error": "ticket service unavailable"})
			return
		}
		defer ticketResp.Body.Close()

		var ticketStats struct {
			SoldTickets  int `json:"sold_tickets"`
			ActiveBuyers int `json:"active_buyers"`
		}
		json.NewDecoder(ticketResp.Body).Decode(&ticketStats)

		// fetch waiting-room stats
		waitResp, err := http.Get(waitingRoomURL + "/stats")
		if err != nil {
			c.JSON(500, gin.H{"error": "waiting room unavailable"})
			return
		}
		defer waitResp.Body.Close()

		var waitingStats struct {
			WaitingUsers int `json:"waiting_users"`
		}
		json.NewDecoder(waitResp.Body).Decode(&waitingStats)

		// merge and send clean response
		c.JSON(200, gin.H{
			"sold_tickets":  ticketStats.SoldTickets,
			"active_buyers": ticketStats.ActiveBuyers,
			"waiting_users": waitingStats.WaitingUsers,
		})
	})

	r.Run(":8083")
}
