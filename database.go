package main

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB

func (db *FireDB) Connect() error {
	// Find home directory.
	home, err := os.Getwd()
	if err != nil {
		return err
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(home + "Your JSON File Location")
	config := &firebase.Config{DatabaseURL: "Your Database URL"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
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

	ref := fireDB.NewRef("/path/to/data")
	err = ref.Set(context.Background(), "some value")
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
