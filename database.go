package main

import (
	"context"
	"fmt"

	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// create ref at path students/:userId
	ref := client.NewRef("students/" + fmt.Sprint(1))

	if err := ref.Set(context.TODO(), map[string]interface{}{
		"name":       "Jane Doe",
		"age":        "7",
		"class":      "Tech Explorer",
		"instructor": "Scott Smith"}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Student added/updated successfully!")
}
