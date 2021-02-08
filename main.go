package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	httpClient     = &http.Client{Timeout: 10 * time.Second}
	nameToNickname = map[string]string{
		"alik":  "alikklimenkov",
		"lesha": "shinskix",
		"sasha": "Jakser",
		"vetal": "closeup24",
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("unable to read .env file. %v\n", err)
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	fortniteTrackerApiKey := os.Getenv("FORTNITE_TRACKER_API_KEY")
	webhookUrl := os.Getenv("WEBHOOK_URL")
	port := os.Getenv("PORT")

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookUrl))
	if err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":"+port, nil)
	fmt.Println("start listen :" + port)

	client := FortniteTrackerClient{
		ApiKey:  fortniteTrackerApiKey,
		BaseUrl: "https://api.fortnitetracker.com/v1",
	}

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			chatID := update.Message.Chat.ID
			command := update.Message.Command()
			switch command {
			case "start":
				continue
			case "stats":
				sendPlayerStats(chatID, client, bot, update.Message.CommandArguments())
			case "alik", "vetal", "lesha", "sasha":
				sendPlayerStats(chatID, client, bot, nameToNickname[command])
			case "team":
				var nicknames []string
				for _, value := range nameToNickname {
					nicknames = append(nicknames, value)
				}
				sendPlayerStats(chatID, client, bot, nicknames...)
			default:
				sendUnknownCommand(chatID, bot)
			}
		}
	}
}

func sendPlayerStats(chatID int64, client FortniteTrackerClient, bot *tgbotapi.BotAPI, nicknames ...string) {
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
