package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var HackTeamCollection *mongo.Collection
var HackMemberCollection *mongo.Collection

func InitDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(context.Background())

	fmt.Println("Connected to MongoDB!")
	HackTeamCollection = client.Database("Hackathon").Collection("HackTeam")
	HackMemberCollection = client.Database("Hackathon").Collection("HackMember")

}
