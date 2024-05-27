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
