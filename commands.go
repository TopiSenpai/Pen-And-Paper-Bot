package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink/api"
)

func onSoundAddCommand(event *core.SlashCommandEvent) {
	_ = event.DeferCreate(true)

	soundName := event.Options["name"].String()

	query := event.Options["source"].String()

	// TODO: ytsearch: if no link or so

	dgolink.RestClient().LoadItemHandler(query, api.NewResultHandler(
		func(track api.Track) {
			tr
		},
		func(playlist *api.Playlist) {

		},
		func(tracks []api.Track) {

		},
		func() {
			event
		},
		func(e *api.Exception) {

		}),
	)
}

func onSoundRemoveCommand(event *core.SlashCommandEvent) {
	soundName := event.Options["name"].String()
}

func onSoundListCommand(event *core.SlashCommandEvent) {

}