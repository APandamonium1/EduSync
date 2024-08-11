package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Use godot package to load/read the .env file and
// return the value of the key (for local env)
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// InitializeFirebase initializes the Firebase app and sets the global firebaseClient variable
func initializeFirebase() error {
	ctx := context.Background()

	databaseURL, found := os.LookupEnv("DATABASE_URL")
	if !found {
		log.Fatalf("DATABASE_URL is not set in the environment variables")
	}
	opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")

	// databaseURL := goDotEnvVariable("DATABASE_URL")
	// if databaseURL == "" {
	// 	return fmt.Errorf("DATABASE_URL is not set in the environment variables")
	// }
	// opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	conf := &firebase.Config{DatabaseURL: databaseURL}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	var firebaseApp = app

	err = initDB(firebaseApp)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	return nil
}
