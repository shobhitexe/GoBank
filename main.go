package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	store, err := NewPostrgreStore()

	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8080", store)
	server.Run()

}
