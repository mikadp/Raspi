//Here make db connection and querys

import (
	"fmt"
	"log"
)

func dbConnection() {
	db, err := sql.Open("")
	if err != nil {
		log.Fatalf("Cannot create the connection: %s", err)
		return nil, err
	}
}
