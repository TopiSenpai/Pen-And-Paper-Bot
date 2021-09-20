package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgolink"
)

type Sound struct {
	Name    string
	Track64 string
}

func (s Sound) ToTrack() disgolink.Track {
	return &disgolink.DefaultTrack{
		Base64Track: &s.Track64,
	}
}

func addNewSound(event *core.SlashCommandEvent, name string, track disgolink.Track) {
	if _, ok := sounds[event.User.ID]; !ok {
		sounds[event.User.ID] = map[string]string{}
	}

	sounds[event.User.ID][name] = track.Track()
	_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("added new sound `%s` with [source](%s)", name, track.Info().URI()).Build())
}
