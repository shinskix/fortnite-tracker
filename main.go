package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"

	"io"
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
	defaultRowHeader = []string{"Mode", "Wins", "WinRate(%)", "Kills", "KD", "Rating"}
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
				playerInfo, err := client.PlayerInfo(PC, update.Message.CommandArguments())
				if err != nil {
					log.Println(err)
				} else {
					bot.Send(prepareStats(chatID, playerInfo))
				}
			case "alik", "vetal", "lesha", "sasha":
				if nickname, exists := nameToNickname[command]; exists {
					playerInfo, err := client.PlayerInfo(PC, nickname)
					if err != nil {
						log.Println(err)
					} else {
						bot.Send(prepareStats(chatID, playerInfo))
					}
				}
			case "team":
				group := PlayerInfoGroup{}
				for _, nickname := range nameToNickname {
					info, err := client.PlayerInfo(PC, nickname)
					if err != nil {
						log.Println(err)
					} else {
						group.Players = append(group.Players, *info)
					}
					time.Sleep(2 * time.Second)
				}
				bot.Send(prepareStats(chatID, &group))
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Unknown or not yet implemented command."))
			}
		}
	}
}

func prepareStats(chatID int64, asciiStats AsciiTransformable) tgbotapi.MessageConfig {
	buf := new(bytes.Buffer)
	asciiStats.transform(buf)
	msg := tgbotapi.NewMessage(chatID, "<pre>"+buf.String()+"</pre>")
	msg.ParseMode = "html"
	return msg
}

func statsToRow(modeStats GameModeStats) []string {
	return []string{
		modeStats.Wins.DisplayValue,
		modeStats.WinRatio.DisplayValue,
		modeStats.Kills.DisplayValue,
		modeStats.KD.DisplayValue,
		modeStats.TrnRating.DisplayValue,
	}
}

type AsciiTransformable interface {
	transform(out io.Writer)
}

func (player *PlayerInfo) transform(out io.Writer) {
	data := [][]string{
		append([]string{"solo"}, statsToRow(player.Stats.Solo)...),
		append([]string{"duos"}, statsToRow(player.Stats.Duos)...),
		append([]string{"squads"}, statsToRow(player.Stats.Squads)...),
	}
	table := tablewriter.NewWriter(out)
	table.SetHeader(defaultRowHeader)
	table.SetFooter([]string{"", "", "", "", "Player", player.Name})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}

type PlayerInfoGroup struct {
	Players []PlayerInfo
}

func (group *PlayerInfoGroup) transform(out io.Writer) {
	table := tablewriter.NewWriter(out)
	table.SetHeader(append([]string{"Nickname"}, defaultRowHeader...))

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "solo"}, statsToRow(player.Stats.Solo)...))
	}

	table.Append([]string{"", "", "", "", "", "", ""})

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "duos"}, statsToRow(player.Stats.Duos)...))
	}

	table.Append([]string{"", "", "", "", "", "", ""})

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "squads"}, statsToRow(player.Stats.Squads)...))
	}

	table.SetBorder(false)
	table.Render()
}
