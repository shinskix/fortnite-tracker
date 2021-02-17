package main

import (
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)
import "context"

func main() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://fortnit-elves-bot-default-rtdb.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("./firebase-private-key.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	info := &PlayerInfo{}
	fd, err := os.Open("./stats.json")
	if err != nil {
		log.Fatal(err)
	}
	stat, err := fd.Stat()
	if err != nil {
		log.Fatal(err)
	}
	bytes := make([]byte, stat.Size())
	fd.Read(bytes)
	json.Unmarshal(bytes, info)
	if _, _, err = client.Collection("stats").Add(ctx, info); err != nil {
		log.Fatal(err)
	}
}
