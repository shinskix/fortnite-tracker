package main

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

var (
	defaultRowHeader = []string{"Mode", "Wins", "WinRate(%)", "Kills", "KD", "Rating"}
)

type AsciiTransformable interface {
	Transform(out io.Writer)
}

func (player *PlayerInfo) Transform(out io.Writer) {
	data := [][]string{
		append([]string{"solo"}, statsToRow(player.Stats.Solo)...),
		append([]string{"duos"}, statsToRow(player.Stats.Duos)...),
		append([]string{"squads"}, statsToRow(player.Stats.Squads)...),
	}
	table := tablewriter.NewWriter(out)
	table.SetHeader(defaultRowHeader)
	table.SetFooter([]string{"", "", "", "", "Player", player.Name})
	table.AppendBulk(data)
	table.Render()
}

func (group *PlayerInfoGroup) Transform(out io.Writer) {
	rowSeparator := []string{"", "", "", "", "", "", ""}

	table := tablewriter.NewWriter(out)
	table.SetHeader(append([]string{"Nickname"}, defaultRowHeader...))

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "solo"}, statsToRow(player.Stats.Solo)...))
	}

	table.Append(rowSeparator)

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "duos"}, statsToRow(player.Stats.Duos)...))
	}

	table.Append(rowSeparator)

	for _, player := range group.Players {
		table.Append(append([]string{player.Name, "squads"}, statsToRow(player.Stats.Squads)...))
	}

	table.Render()
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
