package main

import (
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var redisClient *redis.Client
var db *sql.DB

const (
	redisCacheKeyPrefix = "promotion:"
	cacheSize           = 3 // just for testing LRU, this has to be a big number
	cacheExpiration     = 1 * time.Hour
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp(database:3306)/"+os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Println("Error connecting to mysql:", err)
		return
	}

	err = handleCSV(db)
	if err != nil {
		fmt.Println("Handle csv error:", err)
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})

	setupRoutes()
}
