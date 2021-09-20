module github.com/TopiSenpai/SoundBoard-Bot

go 1.16

replace (
	github.com/DisgoOrg/disgolink => ../disgolink
)

require (
	github.com/DisgoOrg/disgo v0.5.12-0.20210918181853-f7cd2b5ba91e
	github.com/DisgoOrg/disgolink v0.2.5-0.20210919155345-1f0af3080486
	github.com/sirupsen/logrus v1.8.1
)
