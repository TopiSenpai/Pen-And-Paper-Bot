package main

import (
	"github.com/DisgoOrg/disgo/core"
)

func onPlayButton(event *core.ButtonClickEvent, soundName string) {
	for _, sound := range Sounds {
		if sound.Name == soundName {
			guildID := *event.Interaction.GuildID
			if event.Bot().Caches.VoiceStateCache().Get(guildID, event.Bot().ApplicationID) == nil {
				voiceState := event.Bot().Caches.VoiceStateCache().Get(guildID, event.Interaction.Member.User.ID)
				if voiceState == nil {
					_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContent("please connect to a voice channel").Build())
					return
				}
				_ = voiceState.Channel().Connect()
			}
			dgolink.Player(*event.GuildID).Play(sound.ToTrack())

			_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContentf("playing %s", sound.Name).Build())
			return
		}
	}
}

func onStopButton(event *core.ButtonClickEvent) {
	dgolink.Player(*event.GuildID).Stop()
	_ = event.Create(core.NewMessageCreateBuilder().
		SetEphemeral(true).
		SetContent("stopped sounds").
		Build(),
	)
}
