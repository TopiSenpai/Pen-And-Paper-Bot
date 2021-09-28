package main

import (
	"github.com/DisgoOrg/disgo/core"
	"math"
	"strconv"
)

func onPlayButton(event *core.ButtonClickEvent, soundName string) {
	sound := getUserSound(event.User.ID, soundName)
	if sound == nil {
		return
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
	_ = dgolink.Player(*event.GuildID).Play(sound.ToTrack())

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

func onPageButton(event *core.ButtonClickEvent, pageStr string) {
	page, _ := strconv.Atoi(pageStr)
	userSounds, ok := sounds[event.User.ID]
	if !ok || len(userSounds) == 0 {
		_ = event.Create(core.NewMessageCreateBuilder().SetContent("you have no sounds added to your soundboard").SetEphemeral(true).Build())
		return
	}

	// 20 sounds per page
	pages := int(math.Ceil(float64(len(userSounds)) / float64(20)))

	builder := core.NewMessageUpdateBuilder().SetContentf("page %d/%d", page + 1, pages).SetActionRows(buildActionRows(page, pages, userSounds)...)
	_ = event.Update(builder.Build())
}
