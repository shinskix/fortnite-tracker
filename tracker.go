package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Platform string

const (
	PC          Platform = "pc"
	PLAYSTATION Platform = "psn"
	XBOX        Platform = "xbl"
)

type UserInfo struct {
	UserId    string        `json:"accountId"`
	UserName  string        `json:"epicUserHandle"`
	UserStats FortniteStats `json:"stats"`
}

type FortniteStats struct {
	Solo   GameModeStats `json:"p2"`
	Duos   GameModeStats `json:"p10"`
	Squads GameModeStats `json:"p9"`
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
	ApiKey  string
	BaseUrl string
}

type Request struct {
	FTNClient *FortniteTrackerClient
	Method    string
	URL       string
}

func (req *Request) execute() ([]byte, error) {
	request, err := http.NewRequest(req.Method, req.FTNClient.BaseUrl+req.URL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("TRN-Api-Key", req.FTNClient.ApiKey)
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (client *FortniteTrackerClient) PlayerInfo(platform Platform, nickname string) (*UserInfo, error) {
	req := Request{
		client,
		"GET",
		fmt.Sprintf("/profile/%s/%s", platform, nickname),
	}
	userInfoJson, err := req.execute()
	if err != nil {
		return nil, err
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal(userInfoJson, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
