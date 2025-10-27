package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe represents a single recipe document in MongoDB.
type Recipe struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Title       string   `json:"title" bson:"title"`
	Description string   `json:"desc" bson:"desc"`
	Categories  []string `json:"categories" bson:"categories"`
	Ingredients []string `json:"ingredients" bson:"ingredients"`
	Directions  []string `json:"directions" bson:"directions"`
	Calories    float64  `json:"calories" bson:"calories"`
	Protein     float64  `json:"protein" bson:"protein"`
	Fat         float64  `json:"fat" bson:"fat"`
	Sodium      float64  `json:"sodium" bson:"sodium"`
	Rating      float64  `json:"rating" bson:"rating"`
	Date        string   `json:"date" bson:"date"`
}
