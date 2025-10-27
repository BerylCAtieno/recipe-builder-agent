package main

import (
	"log"
	"os"

	// Replace with your actual project path
	"github.com/BerylCAtieno/recipe-agent/importer"
)

func main() {
	// 1. Get Environment Variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("Error: MONGO_URI environment variable not set. Please set it to your Atlas connection string.")
	}

	// 2. Define the path to your JSON file
	jsonFilePath := "data/full_format_recipes.json"

	// 3. Execute the import function
	log.Printf("Starting recipe import from %s...", jsonFilePath)
	if err := importer.ImportRecipes(mongoURI, jsonFilePath); err != nil {
		log.Fatalf("Data import failed: %v", err)
	}
	log.Println("--- Data Import COMPLETE ---")
}
