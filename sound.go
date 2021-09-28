package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink"
)

var sounds = map[discord.Snowflake][]Sound{}

func getUserSound(userID discord.Snowflake, name string) *Sound {
	userSounds, ok := sounds[userID]
	if ok {
		for _, sound := range userSounds {
			if sound.Name == name {
				return &sound
			}
		}
	}
	return nil
}

type Sound struct {
	Name        string `json:"name"`
	TrackBase64 string `json:"track_base_64"`
}

func (s Sound) ToTrack() disgolink.Track {
	return &disgolink.DefaultTrack{
		Base64Track: &s.TrackBase64,
	}
}

func addNewSound(event *core.SlashCommandEvent, name string, track disgolink.Track) {
	if _, ok := sounds[event.User.ID]; !ok {
		sounds[event.User.ID] = []Sound{}
	}

	sounds[event.User.ID] = append(sounds[event.User.ID], Sound{
		Name:        name,
		TrackBase64: track.Track(),
	})
	_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("added new sound `%s` with [source](%s)", name, track.Info().URI()).Build())
}
