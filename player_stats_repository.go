package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

type PlayerStatsRepository struct {
	DatabaseURL string
	Opts        []option.ClientOption
}

var (
	databaseClient  *db.Client
	databaseContext = context.Background()
)

func (r *PlayerStatsRepository) initialize() error {
	conf := &firebase.Config{DatabaseURL: r.DatabaseURL}
	app, err := firebase.NewApp(databaseContext, conf, r.Opts...)
	if err != nil {
		return err
	}
	databaseClient, err = app.Database(databaseContext)
	if err != nil {
		return err
	}
	return nil
}

func (r *PlayerStatsRepository) storeFortniteStats(info *PlayerInfo) error {
	ref := databaseClient.NewRef("fortnite/stats")
	playerRef := ref.Child(info.Name)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	timeRef := playerRef.Child(timestamp)
	return timeRef.Set(databaseContext, info)
}
