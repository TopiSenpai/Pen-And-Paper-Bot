package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink"
	"regexp"
)

var URLPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

func onSoundAddCommand(event *core.SlashCommandEvent) {
	err := event.DeferCreate(true)
	if err != nil {
		event.Bot().Logger.Error("failed to DeferCreate(true) interaction: ", err)
	}

	name := event.Options["name"].String()
	source := event.Options["source"].String()
	query := source

	if !URLPattern.MatchString(query) {
		query = disgolink.SearchTypeYoutube.Apply(query)
	}

	dgolink.RestClient().LoadItemHandler(query, disgolink.NewResultHandler(
		func(track disgolink.Track) {
			addNewSound(event, name, track)
		},
		func(playlist *disgolink.Playlist) {
			track := playlist.SelectedTrack()
			if track == nil {
				track = playlist.Tracks[0]
			}
			addNewSound(event, name, track)
		},
		func(tracks []disgolink.Track) {
			addNewSound(event, name, tracks[0])
		},
		func() {
			_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("no track for `%s` found", source).Build())
		},
		func(e *disgolink.Exception) {
			_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("error while getting track for `%s`", source).Build())
		}),
	)
}

func onSoundRemoveCommand(event *core.SlashCommandEvent) {
	//soundName := event.Options["name"].String()
}

func onSoundListCommand(event *core.SlashCommandEvent) {

}

func onSoundBoardCommand(event *core.SlashCommandEvent) {
	userSounds, ok := sounds[event.User.ID]
	if !ok || len(userSounds) == 0 {
		_ = event.Create(core.NewMessageCreateBuilder().SetContent("you have no sounds added to your soundboard").SetEphemeral(true).Build())
		return
	}

	builder := core.NewMessageCreateBuilder().SetContent("hit a button to play your sounds!").SetEphemeral(true)
	actionRow := core.NewActionRow()
	for name, _ := range userSounds {
		if len(actionRow.Components) == 5 {
			builder.AddActionRows(actionRow)
			actionRow = core.NewActionRow()
			if len(builder.Components) == 4 {
				break
			}
		}
		actionRow = actionRow.AddComponents(core.NewPrimaryButton(name, "play:"+name, nil))
	}
	if len(actionRow.Components) > 0 {
		builder.AddActionRows(actionRow)
	}

	builder.AddActionRow(core.NewSecondaryButton("pause", "pause", nil), core.NewDangerButton("stop", "stop", nil))
	event.Create(builder.Build())
}
