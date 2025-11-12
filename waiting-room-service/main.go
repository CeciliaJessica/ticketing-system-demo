package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis.ticketing.svc.cluster.local:6379",
	})
)

const MAX_ACTIVE_BUYERS = 500

func main() {
	ticketServiceURL := os.Getenv("TICKET_SERVICE_URL")

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	r := gin.Default()

	r.GET("/enter", func(c *gin.Context) {
		active, err := rdb.Get(ctx, "active_buyers").Int()
		if err == redis.Nil {
			active = 0
		} else if err != nil {
			c.JSON(500, gin.H{"error": "redis read failed"})
			return
		}

		if active >= MAX_ACTIVE_BUYERS {
			// too many active buyers
			if err := rdb.Incr(ctx, "waiting_users").Err(); err != nil {
				c.JSON(500, gin.H{"error": "failed to track waiting user"})
				return
			}
			defer rdb.Decr(ctx, "waiting_users")

			c.JSON(429, gin.H{
				"status":  "waiting",
				"message": "System busy — please retry soon.",
			})
			return
		}

		// forward to ticket service
		resp, err := http.Get(ticketServiceURL + "/tickets")
		if err != nil {
			c.JSON(500, gin.H{"error": "could not reach ticket service"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
	})

	r.GET("/stats", func(c *gin.Context) {
		waiting, _ := rdb.Get(ctx, "waiting_users").Int()
		c.JSON(200, gin.H{"waiting_users": waiting})
	})

	r.Run(":8082")
}
