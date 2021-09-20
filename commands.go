package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink/api"
)

func onSoundAddCommand(event *core.SlashCommandEvent) {
	err := event.DeferCreate(true)
	if err != nil {
		event.Bot().Logger.Error("failed to DeferCreate(true) interaction: ", err)
	}

	name := event.Options["name"].String()
	source := event.Options["source"].String()
	query := source

	// TODO: ytsearch: if no link or so

	dgolink.RestClient().LoadItemHandler(query, api.NewResultHandler(
		func(track api.Track) {
			addNewSound(event, name, track)
		},
		func(playlist *api.Playlist) {
			track := playlist.SelectedTrack()
			if track == nil {
				track = playlist.Tracks[0]
			}
			addNewSound(event, name, track)
		},
		func(tracks []api.Track) {
			addNewSound(event, name, tracks[0])
		},
		func() {
			_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("no track for `%s` found", source).Build())
		},
		func(e *api.Exception) {
			_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("error while getting track for `%s`", source).Build())
		}),
	)
}

func addNewSound(event *core.SlashCommandEvent, name string, track api.Track) {
	if _, ok := sounds[*event.GuildID]; !ok {
		sounds[*event.GuildID] = []Sound{}
	}

	sounds[*event.GuildID] = append(sounds[*event.GuildID], Sound{
		Name:    name,
		Track64: track.Track(),
	})
	_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("added new sound `%s` with [source](%s)", name, track.Info().URI()).Build())
}

func onSoundRemoveCommand(event *core.SlashCommandEvent) {
	soundName := event.Options["name"].String()
}

func onSoundListCommand(event *core.SlashCommandEvent) {

}
