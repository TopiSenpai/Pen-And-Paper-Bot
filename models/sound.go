package models

import "github.com/DisgoOrg/disgolink/api"

type Sound struct {
	Name string
	Track64 string
}

func (s Sound) ToTrack() api.Track{
	return &api.DefaultTrack{
		Track_: &s.Track64,
	}
}
