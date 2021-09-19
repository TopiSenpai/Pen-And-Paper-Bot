package main

import (
	"github.com/DisgoOrg/disgo/core"
	"strings"
)

func onSlashCommand(event *core.SlashCommandEvent) {
	switch event.CommandName {
	case "sounds":

	}
}

func onButtonClick(event *core.ButtonClickEvent) {
	action := strings.Split(event.CustomID, ":")
	switch action[0] {
	case "play":

	case "stop":

	}
}
