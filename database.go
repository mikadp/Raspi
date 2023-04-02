// Here make db connection and querys
package raspi

import (
	"database/sql"
	"log"
)

func dbConnection() {
	db, err := sql.Open("")
	if err != nil {
		log.Fatalf("Cannot create the connection: %s", err)
		return nil, err
	}
}
