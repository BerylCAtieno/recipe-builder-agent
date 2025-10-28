package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BerylCAtieno/recipe-agent/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Constants for MongoDB connection
const (
	DatabaseName   = "recipes_db"
	CollectionName = "recipes"
	BatchSize      = 500
)

// ImportRecipes handles the connection, JSON reading, and data insertion into mongodb.
func ImportRecipes(mongoURI string, jsonFilePath string) error {
	// Setup MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify the connection
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	log.Println("Successfully connected and pinged MongoDB Atlas.")

	// Read and Unmarshal JSON File
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", jsonFilePath, err)
	}

	var recipes []models.Recipe
	if err := json.Unmarshal(data, &recipes); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	log.Printf("Successfully read and unmarshaled %d recipes from JSON file.", len(recipes))

	// Prepare for Bulk Insertion
	var documents []interface{}
	for _, recipe := range recipes {
		// Note: We cast each recipe struct to the interface{} type
		documents = append(documents, recipe)
	}

	// Get the collection handle
	collection := client.Database(DatabaseName).Collection(CollectionName)
	totalInserted := 0

	// Perform Batch Insertion
	for i := 0; i < len(recipes); i += BatchSize {
		end := i + BatchSize
		if end > len(recipes) {
			end = len(recipes)
		}

		// Get the current batch slice
		batch := recipes[i:end]

		// Convert the batch slice to a slice of interfaces{} for InsertMany
		var documents []interface{}
		for _, recipe := range batch {
			documents = append(documents, recipe)
		}

		// Execute InsertMany for the batch
		// Set a shorter context for the insert itself (e.g., 10 seconds per batch)
		insertCtx, insertCancel := context.WithTimeout(context.Background(), 10*time.Second)
		res, err := collection.InsertMany(insertCtx, documents)
		insertCancel() // Clean up the batch context

		if err != nil {
			return fmt.Errorf("failed to insert batch %d-%d: %w", i, end, err)
		}

		insertedCount := len(res.InsertedIDs)
		totalInserted += insertedCount
		log.Printf("Inserted batch %d-%d. Count: %d. Total inserted: %d", i, end, insertedCount, totalInserted)
	}

	log.Printf("Successfully inserted a total of %d documents into the %s collection.", totalInserted, CollectionName)
	return nil
}
