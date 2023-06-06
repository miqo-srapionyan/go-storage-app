package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Promotion struct {
	UUID           string  `json:"uuid"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}

// API
func setupRoutes() {
	router := gin.Default()
	router.GET("/promotions/:id", getPromotionByID)
	router.Run(":1321")
}

// Find promotion by id, we must use Redis for caching frequently requested promotions
// Use database sharding for performance, and replication to minimizes downtime
func getPromotionByID(context *gin.Context) {
	id := context.Param("id")

	// Check the cache first
	cachedPromotion, err := getPromotionFromCache(id, redisClient)
	if err == nil {
		context.JSON(http.StatusOK, cachedPromotion)
		return
	}

	// Cache miss, fetch from the database
	promotion, err := getPromotionFromDatabase(id)
	if err == sql.ErrNoRows {
		context.JSON(http.StatusNotFound, gin.H{"error": "Promotion not found"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Add the fetched promotion to cache
	err = addPromotionToCache(id, promotion)
	if err != nil {
		log.Println("Failed to add promotion to cache:", err)
	}

	context.JSON(http.StatusOK, promotion)
}
