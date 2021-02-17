package main

import (
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
	"time"
)
import "context"

func mainTest() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://fortnit-elves-bot-default-rtdb.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("./firebase-private-key.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Database(ctx)
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
	ref := client.NewRef("fortnite/stats")
	playerRef := ref.Child(info.Name)
	timeRef := playerRef.Child(strconv.FormatInt(time.Now().Unix(), 10))
	newStat, err := timeRef.Push(ctx, nil)
	if err != nil {
		log.Fatalln("error pushing stat node:", err)
	}
	if err = newStat.Set(ctx, info); err != nil {
		log.Fatal(err)
	}
}
