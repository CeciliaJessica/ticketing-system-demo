package main

import (
	"context"
	"net/http"
	"os"
	"time"

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

	httpClient := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	r := gin.Default()

	r.GET("/enter", func(c *gin.Context) {
		// Reserve a slot first
		active, err := rdb.Incr(ctx, "active_buyers").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": "redis incr failed"})
			return
		}

		// If over capacity -> rollback and 429
		if active > MAX_ACTIVE_BUYERS {
			_ = rdb.Decr(ctx, "active_buyers").Err()
			_ = rdb.Incr(ctx, "waiting_users").Err()
			defer func() { _ = rdb.Decr(ctx, "waiting_users").Err() }()

			c.JSON(429, gin.H{
				"status":  "waiting",
				"message": "System busy — please retry soon.",
			})
			return
		}

		// Always release slot when request finishes
		defer func() { _ = rdb.Decr(ctx, "active_buyers").Err() }()

		// Forward to ticket service
		resp, err := httpClient.Get(ticketServiceURL + "/tickets")
		if err != nil {
			c.JSON(500, gin.H{"error": "could not reach ticket service"})
			return
		}
		defer resp.Body.Close()

		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			ct = "application/json"
		}
		c.DataFromReader(resp.StatusCode, resp.ContentLength, ct, resp.Body, nil)
	})

	r.GET("/stats", func(c *gin.Context) {
		waiting, _ := rdb.Get(ctx, "waiting_users").Int()
		active, _ := rdb.Get(ctx, "active_buyers").Int()
		c.JSON(200, gin.H{"waiting_users": waiting, "active_buyers": active})
	})

	r.Run(":8082")
}
