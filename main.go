package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

type Config struct {
	PlayerNameToNickname    map[string]string `envconfig:"FORTNITE_PLAYERS" required:"true"`
	FortniteTrackerApiKey   string            `envconfig:"FORTNITE_TRACKER_API_KEY" required:"true"`
	TelegramBotToken        string            `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	WebhookURL              string            `envconfig:"TELEGRAM_WEBHOOK_URL"`
	FirebaseDatabaseURL     string            `envconfig:"FIREBASE_DB_URL" required:"true"`
	FirebaseCredentialsFile string            `envconfig:"FIREBASE_DB_CREDENTIALS"`
	DbSyncCronSpec          string            `envconfig:"CRON_SYNC_STATS_SPEC" required:"true"`
	Port                    string            `envconfig:"PORT" default:"8080"`
}

func main() {
	var appConfig Config
	err := envconfig.Process("", &appConfig)
	if err != nil {
		log.Fatalln("Error loading application config", err)
	}

	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramBotToken)
	if err != nil {
		log.Fatalln("Error creating telegram bot instance", err)
	}
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	if appConfig.WebhookURL != "" {
		// for local development
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(appConfig.WebhookURL))
		if err != nil {
			log.Fatalln("Error setting telegram webhook", err)
		}
	}

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+appConfig.Port, nil)
	fmt.Println("start listen :" + appConfig.Port)

	var elvesNicknames []string
	for _, value := range appConfig.PlayerNameToNickname {
		elvesNicknames = append(elvesNicknames, value)
	}

	fortniteTrackerClient := FortniteTrackerClient{
		ApiKey:  appConfig.FortniteTrackerApiKey,
		BaseUrl: "https://api.fortnitetracker.com/v1",
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	var opts []option.ClientOption
	if appConfig.FirebaseCredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(appConfig.FirebaseCredentialsFile))
	}

	playerStatsRepository := &PlayerStatsRepository{
		appConfig.FirebaseDatabaseURL,
		opts,
	}
	err = playerStatsRepository.initialize()
	if err != nil {
		log.Fatalln("Error initializing repository", err)
	}

	dbCron := cron.New()
	_, err = dbCron.AddFunc(appConfig.DbSyncCronSpec, func() {
		syncStats(elvesNicknames, fortniteTrackerClient, playerStatsRepository)
	})
	if err != nil {
		log.Fatalln("Error registering cron function", err)
	}
	go dbCron.Start()

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			chatID := update.Message.Chat.ID
			command := update.Message.Command()
			switch command {
			case "start":
				continue
			case "stats":
				sendPlayerPhotoStats(chatID, fortniteTrackerClient, bot, update.Message.CommandArguments())
			case "alik", "vetal", "lesha", "sasha":
				sendPlayerPhotoStats(chatID, fortniteTrackerClient, bot, appConfig.PlayerNameToNickname[command])
			case "team":
				sendPlayerPhotoStats(chatID, fortniteTrackerClient, bot, elvesNicknames...)
			default:
				sendUnknownCommand(chatID, bot)
			}
		}
	}
}

func syncStats(elvesNicknames []string, fortniteTrackerClient FortniteTrackerClient, playerStatsRepository *PlayerStatsRepository) {
	for _, nickname := range elvesNicknames {
		log.Printf("Syncing %s stats", nickname)
		playerInfo, err := fortniteTrackerClient.PlayerInfo(PC, nickname)
		if err == nil {
			err = playerStatsRepository.storeFortniteStats(playerInfo)
			if err != nil {
				log.Printf("Failed to store %s stats. Error: %v\n", nickname, err)
			}
		} else {
			log.Printf("Failed to sync %s stats. Error %v\n", nickname, err)
		}
	}
}

func sendPlayerPhotoStats(chatID int64, client FortniteTrackerClient, bot *tgbotapi.BotAPI, nicknames ...string) {
	var playerInfo AsciiTransformable
	var err error
	if len(nicknames) == 1 {
		playerInfo, err = client.PlayerInfo(PC, nicknames[0])
	} else {
		playerInfo, err = client.PlayerInfoGroup(PC, nicknames)
	}
	if err == nil {
		bot.Send(prepareStats(chatID, playerInfo))
	} else {
		log.Println(err)
	}
}

func sendUnknownCommand(chatID int64, bot *tgbotapi.BotAPI) {
	bot.Send(tgbotapi.NewMessage(chatID, "Unknown or not yet implemented command."))
}

func prepareStats(chatID int64, asciiStats AsciiTransformable) tgbotapi.PhotoConfig {
	textBuf := new(bytes.Buffer)
	asciiStats.Transform(textBuf)
	imgBuf := new(bytes.Buffer)
	err := CreateImage(imgBuf, textBuf.String())
	if err != nil {
		log.Println(err)
	}
	return tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "Stats", Bytes: imgBuf.Bytes()})
}
