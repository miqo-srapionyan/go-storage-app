package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

// Insert chunks of data into database using transaction for consistency
func insertData(db *sql.DB, recordsChan <-chan CSVRecord, wg *sync.WaitGroup) {
	defer wg.Done()

	// Begin the transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error on begin txn ", err)
		panic(err.Error())
	}

	stmt, err := tx.Prepare("INSERT INTO promotions (uuid, price, expiration_date) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		err := tx.Rollback()
		if err != nil {
			return
		}
		return
	}
	defer stmt.Close()

	// Process records from the channel until it's closed
	for record := range recordsChan {
		// Execute the INSERT statement with the record values
		_, err := stmt.Exec(record.id, record.price, record.expiration_date)
		if err != nil {
			fmt.Println(err)
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
	}

	// Commit the transaction if all records processed successfully
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		err := tx.Rollback()
		if err != nil {
			return
		}
	}
}

func getPromotionFromDatabase(id string) (*Promotion, error) {
	var promotion Promotion
	err := db.QueryRow(`
		SELECT uuid, price, expiration_date
		FROM promotions
		WHERE id = ?
	`, id).Scan(&promotion.UUID, &promotion.Price, &promotion.ExpirationDate)

	return &promotion, err
}
