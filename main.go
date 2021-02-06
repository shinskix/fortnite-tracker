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
	httpClient = &http.Client{Timeout: 5 * time.Second}
	nameToNick = map[string]string{
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
				bot.Send(playerStats(chatID, update.Message.CommandArguments(), client))
			case "alik", "vetal", "lesha", "sasha":
				if nickname, exists := nameToNick[command]; exists {
					bot.Send(playerStats(chatID, nickname, client))
				}
			case "team":
				bot.Send(groupStats(chatID, client))
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Unknown or not yet implemented command."))
			}
		}
	}
}

func playerStats(chatID int64, nickname string, client FortniteTrackerClient) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "")
	if len(nickname) > 0 {
		playerInfo, err := client.PlayerInfo(PC, nickname)
		if err != nil {
			log.Printf("failed to get player info %v", err)
			msg.Text = "Try again in a few moments"
		} else {
			buf := new(bytes.Buffer)
			writeStats(buf, playerInfo)
			msg.ParseMode = "html"
			msg.Text = "<pre>" + buf.String() + "</pre>"
		}
	}
	return msg
}

func groupStats(chatID int64, client FortniteTrackerClient) tgbotapi.MessageConfig {
	var stats []*PlayerInfo
	for _, nickname := range nameToNick {
		info, err := client.PlayerInfo(PC, nickname)
		if err != nil {
			log.Println(err)
		} else {
			stats = append(stats, info)
		}
		time.Sleep(2 * time.Second)
	}
	buf := new(bytes.Buffer)
	writeGroupStats(buf, stats)
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

func writeStats(out io.Writer, player *PlayerInfo) {
	//{"Total", "92", "3.0%", "5444", "1.81", ""}, TODO parse lifeTimeStats
	data := [][]string{
		append([]string{"Solo"}, statsToRow(player.Stats.Solo)...),
		append([]string{"Duos"}, statsToRow(player.Stats.Duos)...),
		append([]string{"Squads"}, statsToRow(player.Stats.Squads)...),
	}
	table := tablewriter.NewWriter(out)
	table.SetHeader(defaultRowHeader)
	table.SetFooter([]string{"", "", "", "", "Player", player.Name})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}

func writeGroupStats(out io.Writer, playerGroup []*PlayerInfo) {
	table := tablewriter.NewWriter(out)
	table.SetHeader(append([]string{"Nickname"}, defaultRowHeader...))

	for _, player := range playerGroup {
		table.Append(append([]string{player.Name, "solo"}, statsToRow(player.Stats.Solo)...))
	}

	table.Append([]string{"", "", "", "", "", "", ""})

	for _, player := range playerGroup {
		table.Append(append([]string{player.Name, "duos"}, statsToRow(player.Stats.Duos)...))
	}

	table.Append([]string{"", "", "", "", "", "", ""})

	for _, player := range playerGroup {
		table.Append(append([]string{player.Name, "squads"}, statsToRow(player.Stats.Squads)...))
	}

	table.SetBorder(false)
	table.Render()
}
