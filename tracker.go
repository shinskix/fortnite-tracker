package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Platform string

const (
	PC          Platform = "pc"
	PLAYSTATION Platform = "psn"
	XBOX        Platform = "xbl"
)

type PlayerInfo struct {
	ID    string `json:"accountId"`
	Name  string `json:"epicUserHandle"`
	Stats struct {
		Solo   GameModeStats `json:"p2"`
		Duos   GameModeStats `json:"p10"`
		Squads GameModeStats `json:"p9"`
	} `json:"stats"`
}

type PlayerInfoGroup struct {
	Players []PlayerInfo
}

type GameModeStats struct {
	TrnRating StatValue `json:"trnRating"`
	Score     StatValue `json:"score"`
	KD        StatValue `json:"kd"`
	KPM       StatValue `json:"kpg"`
	SPM       StatValue `json:"scorePerMatch"`
	WinRatio  StatValue `json:"winRatio"`
	Wins      StatValue `json:"top1"`
	Kills     StatValue `json:"kills"`
}

type StatValue struct {
	Label        string  `json:"label"`
	Field        string  `json:"field"`
	Category     string  `json:"category"`
	Value        string  `json:"value"`
	Percentile   float32 `json:"percentile"`
	DisplayValue string  `json:"displayValue"`
}

type FortniteTrackerClient struct {
	ApiKey     string
	BaseUrl    string
	HttpClient *http.Client
}

type Request struct {
	Method string
	URL    string
}

func (client *FortniteTrackerClient) execute(req *Request) ([]byte, error) {
	request, err := http.NewRequest(req.Method, client.BaseUrl+req.URL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("TRN-Api-Key", client.ApiKey)
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (client *FortniteTrackerClient) PlayerInfo(platform Platform, nickname string) (*PlayerInfo, error) {
	req := &Request{
		"GET",
		fmt.Sprintf("/profile/%s/%s", platform, nickname),
	}
	resp, err := client.execute(req)
	if err != nil {
		return nil, err
	}
	playerInfo := &PlayerInfo{}
	err = json.Unmarshal(resp, playerInfo)
	if err != nil {
		return nil, err
	}
	return playerInfo, nil
}

func (client *FortniteTrackerClient) PlayerInfoGroup(platform Platform, nicknames []string) (*PlayerInfoGroup, error) {
	var group = PlayerInfoGroup{}
	for _, nickname := range nicknames {
		info, err := client.PlayerInfo(platform, nickname)
		if err == nil {
			group.Players = append(group.Players, *info)
		} else {
			log.Println(err)
		}
		time.Sleep(2 * time.Second)
	}
	if len(group.Players) == 0 {
		return nil, fmt.Errorf("failed to get player info group")
	}
	return &group, nil
}
