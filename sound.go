package main

import "github.com/DisgoOrg/disgolink/api"

type Sound struct {
	Name    string
	Track64 string
}

func (s Sound) ToTrack() api.Track {
	return &api.DefaultTrack{
		Base64Track: &s.Track64,
	}
}
