package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ticket model
type Ticket struct {
	ID         uint   `gorm:"primaryKey"`
	SeatNumber string `gorm:"uniqueIndex"`
	Status     string // free, reserved, sold
	UserEmail  string
}

var (
	db  *gorm.DB
	sem = make(chan struct{}, 100) // limit 100 concurrent buyers
)

func main() {
	// Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Try connecting up to 10 times
	var database *gorm.DB
	var err error
	for i := 1; i <= 10; i++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Postgres not ready (attempt %d/10): %v", i, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL after retries:", err)
	}
	db = database
	log.Println("Connected to PostgreSQL!")

	// Auto migrate table
	if err := db.AutoMigrate(&Ticket{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Seed seats if empty
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

		db.CreateInBatches(tickets, 1000)
		log.Println("Seeding complete!")
	}

	r := gin.Default()

	r.GET("/tickets", buyTicket)
	r.GET("/active", activeBuyers)

	r.Run(":8081")
}

func buyTicket(c *gin.Context) {
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()

		// simulates load
		for i := 0; i < 1e6; i++ {
			_ = math.Sqrt(float64(i))
		}
		time.Sleep(5 * time.Second)

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
	c.JSON(http.StatusOK, gin.H{"active_users": len(sem)})
}
