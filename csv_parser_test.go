package main

import (
	"os"
	"testing"
	"time"
)

// Database must be mocked and used - not done that part as it requires more time to learn mocking in GO

func TestParseCSVRecord(t *testing.T) {
	record := []string{"123", "12.34", "2023-06-06 10:30:00 -0700 MST"}

	expectedID := "123"
	expectedPrice := 12.34
	expectedDate := time.Date(2023, 6, 6, 10, 30, 0, 0, time.FixedZone("-0700", -7*60*60))

	result, err := parseCSVRecord(record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.id != expectedID || result.price != expectedPrice || !result.expiration_date.Equal(expectedDate) {
		t.Errorf("expected {id: %s, price: %.2f, expiration_date: %s}, got %v",
			expectedID, expectedPrice, expectedDate, result)
	}
}

func TestOpenCSVFile(t *testing.T) {
	_, err := openCSVFile("nonexistent.csv")
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestCreateCSVReader(t *testing.T) {
	file, err := os.Open("promotions.csv")
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := createCSVReader(file)

	if reader.Comma != ',' {
		t.Errorf("expected comma separator ',', got '%c'", reader.Comma)
	}

	if reader.FieldsPerRecord != -1 {
		t.Errorf("expected FieldsPerRecord = -1, got %d", reader.FieldsPerRecord)
	}
}
