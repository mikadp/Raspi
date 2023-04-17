// Here make db connection and querys
package raspi

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func dbConnection() (*sql.DB, error) {
	// Read db username etc from .env
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	name := os.Getenv("MYSQL_DATABASE")

	// Create the database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", user, password, name))
	if err != nil {
		log.Fatalf("Cannot create the connection: %s", err)
		return nil, err
	}

	return db, nil
}
