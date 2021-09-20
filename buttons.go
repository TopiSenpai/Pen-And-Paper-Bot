package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink"
)

func onPlayButton(event *core.ButtonClickEvent, soundName string) {
	userSounds, ok := sounds[event.User.ID]
	if !ok {
		return
	}
	sound, ok := userSounds[soundName]
	if !ok {
	}

	memberVoiceState := event.Member.VoiceState()
	if memberVoiceState == nil {
		_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContent("please connect to a voice channel").Build())
		return
	}
	voiceState := event.Guild().SelfMember().VoiceState()
	if voiceState == nil || voiceState.ChannelID == nil || *voiceState.ChannelID != *memberVoiceState.ChannelID {
		_ = memberVoiceState.Channel().Connect()
	}
	_ = dgolink.Player(*event.GuildID).Play(&disgolink.DefaultTrack{Base64Track: &sound})

	actionRows := event.Message.ActionRows()

	actionRows[len(actionRows)-1] = actionRows[len(actionRows)-1].SetComponent("stop", core.NewDangerButton("stop", "stop", nil))

	_ = event.Update(core.NewMessageUpdateBuilder().SetActionRows(actionRows...).Build())

}

func onStopButton(event *core.ButtonClickEvent) {
	_ = event.Bot().AudioController.Disconnect(*event.GuildID)
	if player := dgolink.ExistingPlayer(*event.GuildID); player != nil {
		_ = player.Destroy()
	}
	_ = event.UpdateButton(event.Button().AsDisabled())
}

func onPauseButton(event *core.ButtonClickEvent) {
	player := dgolink.Player(*event.GuildID)
	pause := !player.Paused()
	_ = player.Pause(pause)

	label := "pause"
	if pause {
		label = "resume"
	}

	_ = event.UpdateButton(event.Button().WithLabel(label))
}
