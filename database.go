package main

import (
	"context"
	"fmt"

	"log"
	"os"

	"github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

type FireDB struct {
	*db.Client
}

var fireDB FireDB

// Connect initializes and connects to the Firebase Realtime Database.
// It uses the provided JSON file containing the service account key to authenticate with Firebase.
// The DatabaseURL is set to the Firebase project's Realtime Database URL.
// The function returns an error if any step fails during the initialization or connection process.
func (db *FireDB) Connect() error {
	// Find home directory.
	home, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create a context for the Firebase app.
	ctx := context.Background()

	// Set up the Firebase app with the provided JSON file containing the service account key.
	opt := option.WithCredentialsFile(home + "edusync-firebase.json")
	dotenv := goDotEnvVariable("DATABASE_URL")
	config := &firebase.Config{DatabaseURL: dotenv}

	app, err := firebase.NewApp(ctx, config, opt)
	// opt := option.WithCredentialsFile("edusync-firebase.json")
	// app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	// Initialize the Firebase Realtime Database client.
	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	// Assign the Firebase Realtime Database client to the FireDB struct.
	db.Client = client
	return nil
}

func FirebaseDB() *FireDB {
	return &fireDB
}

func connectToFirebase() error {
	fireDB := FirebaseDB()
	err := fireDB.Connect()
	if err != nil {
		return err
	}

	ref := fireDB.NewRef("/")
	err = ref.Set(context.Background(), map[string]string{
		"name":       "Jane Doe",
		"age":        "7",
		"class":      "Tech Explorer",
		"instructor": "Scott Smith",
	})
	if err != nil {
		return err
	}
	return nil
}

// CRUD operations
func (db *FireDB) Create(refPath string, data interface{}) error {
	ref := db.NewRef(refPath)
	return ref.Set(context.Background(), data)
}

func (db *FireDB) Read(refPath string, dest interface{}) error {
	ref := db.NewRef(refPath)
	return ref.Get(context.Background(), dest)
}

func (db *FireDB) Update(refPath string, data map[string]interface{}) error {
	ref := db.NewRef(refPath)
	return ref.Update(context.Background(), data)
}

func (db *FireDB) Delete(refPath string) error {
	ref := db.NewRef(refPath)
	return ref.Delete(context.Background())
}

//todo: set up firebase, finish connecting to database
