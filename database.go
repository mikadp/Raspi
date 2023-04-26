// Here make db connection and querys
package raspi

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func setupDatabaseConnection() {
	// Read db username etc from .env
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	name := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("ENV MYSQL_PORT")

	// Create the database connection
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name))
	if err != nil {
		log.Fatalf("Cannot create the connection: %s", err)
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
