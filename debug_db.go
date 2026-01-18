package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Collection struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	path := "backend/models.db"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Database not found at %s", path)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var tables []string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
	fmt.Println("Tables:", tables)

	var collections []Collection
	// SELECT * to see what columns exist effectively
	if err := db.Find(&collections).Error; err != nil {
		fmt.Printf("Error querying collections: %v\n", err)
	} else {
		fmt.Printf("Found %d collections:\n", len(collections))
		for _, c := range collections {
			fmt.Printf("ID: %d, Name: '%s', Desc: '%s', Created: %v\n", c.ID, c.Name, c.Description, c.CreatedAt)
		}
	}
}
