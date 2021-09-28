package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink"
	"math"
	"regexp"
	"strconv"
)

var URLPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")

func onSoundAddCommand(event *core.SlashCommandEvent) {

	name := event.Options["name"].String()

	if sound := getUserSound(event.User.ID, name); sound != nil {
		_ = event.Create(core.NewMessageCreateBuilder().SetContentf("you already have a sound with the name: %s", name).SetEphemeral(true).Build())
		return
	}

	err := event.DeferCreate(true)
	if err != nil {
		event.Bot().Logger.Error("failed to DeferCreate(true) interaction: ", err)
	}

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
	name := event.Options["name"].String()
	if _, ok := sounds[event.User.ID]; !ok {
		_ = event.Create(core.NewMessageCreateBuilder().SetContentf("no sounds found for you").SetEphemeral(true).Build())
		return
	}

	var soundI *int
	for i, sound := range sounds[event.User.ID] {
		if sound.Name == name {
			soundI = &i
			break
		}
	}

	if soundI == nil {
		_ = event.Create(core.NewMessageCreateBuilder().SetContentf("no sound with name: %s found", name).SetEphemeral(true).Build())
		return
	}
	sounds[event.User.ID] = append(sounds[event.User.ID][:*soundI], sounds[event.User.ID][*soundI+1:]...)

	_ = event.Create(core.NewMessageCreateBuilder().SetContentf("removed sound: %s", name).Build())

}

func onSoundBoardCommand(event *core.SlashCommandEvent) {
	userSounds, ok := sounds[event.User.ID]
	if !ok || len(userSounds) == 0 {
		_ = event.Create(core.NewMessageCreateBuilder().SetContent("you have no sounds added to your soundboard").SetEphemeral(true).Build())
		return
	}

	// 20 sounds per page
	pages := int(math.Ceil(float64(len(userSounds)) / float64(20)))

	builder := core.NewMessageCreateBuilder().SetContentf("page 1/%d", pages).SetEphemeral(true).SetActionRows(buildActionRows(0, pages, userSounds)...)
	_ = event.Create(builder.Build())
}

func buildActionRows(page int, pages int, sounds []Sound) []core.ActionRow {
	var rows []core.ActionRow

	actionRow := core.NewActionRow()
	for i := page * 20; i < len(sounds); i++ {
		if len(actionRow.Components) == 5 {
			rows = append(rows, actionRow)
			actionRow = core.NewActionRow()
			if len(rows) == 4 {
				break
			}
		}
		actionRow = actionRow.AddComponents(core.NewPrimaryButton(sounds[i].Name, "play:"+sounds[i].Name, nil))
	}
	if len(actionRow.Components) > 0 {
		rows = append(rows, actionRow)
	}

	actionRow = core.NewActionRow()
	if pages > 1 {
		left := core.NewSecondaryButton("left", "page:"+strconv.Itoa(page-1), &discord.Emoji{Name: "⬅"})
		if page == 0 {
			left = left.AsDisabled()
		}
		right := core.NewSecondaryButton("right", "page:"+strconv.Itoa(page+1), &discord.Emoji{Name: "➡"})
		if page == pages-1 {
			right = right.AsDisabled()
		}
		actionRow = actionRow.AddComponents(
			left,
			right,
		)
	}

	rows = append(rows, actionRow.AddComponents(core.NewSecondaryButton("refresh", "page:"+strconv.Itoa(page), nil), core.NewSecondaryButton("pause", "pause", nil), core.NewDangerButton("stop", "stop", nil)))
	return rows
}
