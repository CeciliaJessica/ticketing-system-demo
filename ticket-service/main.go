package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Ticket struct {
	ID         uint64 `gorm:"primaryKey"`
	SeatNumber string `gorm:"uniqueIndex"`
	Status     string
	UserEmail  string
}

var (
	db  *gorm.DB
	sem = make(chan struct{}, 100)
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "redis.ticketing.svc.cluster.local:6379",
})

func main() {
	// PostgreSQL connection
	cockroachURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:26257/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var database *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		database, err = gorm.Open(postgres.Open(cockroachURL), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("CockroachDB not ready (attempt %d/10): %v", i, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to CockroachDB after retries:", err)
	}

	db = database
	log.Println("Connected to CockroachDB!")

	if err := db.AutoMigrate(&Ticket{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Seed tickets if empty
	var count int64
	db.Model(&Ticket{}).Count(&count)
	if count == 0 {
		rows := 50
		seatsPerRow := 1000
		totalSeats := rows * seatsPerRow

		log.Printf("Seeding %d seats...", totalSeats)
		var tickets []Ticket
		for r := 0; r < rows; r++ {
			rowLetter := string('A' + r)
			for s := 1; s <= seatsPerRow; s++ {
				seatNumber := fmt.Sprintf("%d%s", s, rowLetter)
				tickets = append(tickets, Ticket{SeatNumber: seatNumber, Status: "free"})
			}
		}
		db.CreateInBatches(tickets, 500)
		log.Println("Seeding complete!")
	}

	r := gin.Default()
	r.GET("/tickets", buyTicket)
	r.GET("/active", activeBuyers)
	r.GET("/stats", getStats)

	r.Run(":8081")
}

func buyTicket(c *gin.Context) {
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()

		ctx := context.Background()

		// increment active buyers counter
		if err := rdb.Incr(ctx, "active_buyers").Err(); err != nil {
			log.Printf("Redis INCR failed: %v", err)
		}
		// ensure decrement always runs
		defer func() {
			if err := rdb.Decr(ctx, "active_buyers").Err(); err != nil {
				log.Printf("Redis DECR failed: %v", err)
			}
		}()

		// simulate heavy computation / delay
		for i := 0; i < 1e6; i++ {
			_ = math.Sqrt(float64(i))
		}
		time.Sleep(5 * time.Second)

		// try to reserve a free ticket
		var ticket Ticket
		tx := db.Where("status = ?", "free").First(&ticket)
		if tx.Error != nil {
			c.JSON(http.StatusGone, gin.H{"message": "All tickets sold out!"})
			return
		}

		ticket.Status = "sold"
		db.Save(&ticket)

		c.JSON(http.StatusOK, gin.H{
			"message": "Ticket purchased successfully",
			"seat":    ticket.SeatNumber,
		})

	default:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": "Too many users buying tickets right now, please wait.",
		})
	}
}

func activeBuyers(c *gin.Context) {
	ctx := context.Background()
	active, err := rdb.Get(ctx, "active_buyers").Int()
	if err == redis.Nil {
		active = 0
	} else if err != nil {
		log.Printf("Redis read error: %v", err)
		active = 0
	}
	c.JSON(http.StatusOK, gin.H{"active_users": active})
}

func getStats(c *gin.Context) {
	var soldCount int64
	db.Model(&Ticket{}).Where("status = ?", "sold").Count(&soldCount)

	ctx := context.Background()
	active, err := rdb.Get(ctx, "active_buyers").Int()
	if err == redis.Nil {
		active = 0
	} else if err != nil {
		log.Printf("Redis read error: %v", err)
		active = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"sold_tickets":  soldCount,
		"active_buyers": active,
	})
}
