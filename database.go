// Here make db connection and querys
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

func setupDatabaseConnection() {
	// Read db username etc from .env
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// Create the database connection
	var err error
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Cannot create the connection: %s", err)
	}

	//Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to the database: %s", err)
	}
}

func getPhoneNumbers() ([]string, error) {
	// Query the phonenumbers table to retrieve all phone numbers
	rows, err := db.Query("SELECT number FROM phonenumbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the phone numbers
	var phoneNumbers []string
	for rows.Next() {
		var phoneNumber string
		if err := rows.Scan(&phoneNumber); err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}

	return phoneNumbers, nil
}

func getTelegramAPI() (string, error) {
	//query the api
	var botToken string
	err := db.QueryRow("SELECT botToken FROM telegram")
	if err != nil {
		return "", fmt.Errorf("Error loading botToken from database: %v", err)
	}

	return botToken, nil
}
