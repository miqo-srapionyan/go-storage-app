package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	csvFilePath = "promotions.csv"
	numWorkers  = 10 // Number of concurrent workers for processing CSV records
)

// CSVRecord Custom struct representing a record
type CSVRecord struct {
	id              string
	price           float64
	expiration_date time.Time
}

// Handle csv import, this is works only once, and designed for this task
// To be able to upload or receive new file we must have an API endpoint or point our code to that directory
// This is done with goroutines for parallel processing
func handleCSV(db *sql.DB) error {
	file, err := openCSVFile(csvFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer closeCSVFile(file)

	// Create a channel to receive CSV records
	recordsChan := make(chan CSVRecord)

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start multiple workers to process the records concurrently
	for i := 0; i < numWorkers; i++ {
		go insertData(db, recordsChan, &wg)
	}

	// Read the CSV file line by line and send records to the channel
	csvReader := createCSVReader(file)
	startTime := time.Now()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			continue
		}

		// Convert the CSV record to custom CSVRecord struct
		r, _ := parseCSVRecord(record)

		// Send the record to the channel for processing
		recordsChan <- r
	}

	// Close the channel after sending all records
	close(recordsChan)

	// Wait for all workers to finish
	wg.Wait()

	// Calculate and print the total execution time
	elapsedTime := time.Since(startTime)
	fmt.Printf("Import completed in %s\n", elapsedTime)

	return err
}

func parseCSVRecord(record []string) (CSVRecord, error) {
	price, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return CSVRecord{}, fmt.Errorf("error parsing price: %w", err)
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700 MST", record[2])
	if err != nil {
		return CSVRecord{}, fmt.Errorf("error parsing date: %w", err)
	}

	newRecord := CSVRecord{
		id:              record[0],
		price:           price,
		expiration_date: date,
	}
	return newRecord, nil
}

func createCSVReader(file *os.File) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	return reader
}

func openCSVFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func closeCSVFile(file *os.File) {
	file.Close()
}
