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
)

// ImportRecipes handles the connection, JSON reading, and data insertion into mongodb.
func ImportRecipes(mongoURI string, jsonFilePath string) error {
	// Setup MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// Perform Bulk Insertion
	res, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert recipes: %w", err)
	}

	log.Printf("Successfully inserted %d documents into the %s collection.", len(res.InsertedIDs), CollectionName)
	return nil
}
