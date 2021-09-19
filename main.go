package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgolink"
	"github.com/DisgoOrg/disgolink/api"
	"github.com/TopiSenpai/SoundBoard-Bot/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var Sounds = []models.Sound{
	{
		Name:    "Wooosh",
		Track64: "QAAAfgIAFFdob29zaCBTb3VuZCBFZmZlY3RzABAyTWlycm9yc0RpYWxvZ3VlAAAAAAAAv2gAC0hEUlZ6d05rVjIwAAEAK2h0dHBzOi8vd3d3LnlvdXR1YmUuY29tL3dhdGNoP3Y9SERSVnp3TmtWMjAAB3lvdXR1YmUAAAAAAAAAAA==",
	},
}

var token = os.Getenv("token")

//var logWebhookToken = os.Getenv("log_webhook_token")

var logger = logrus.New()
var httpClient = http.DefaultClient
var bot *core.Bot
var dgolink api.Disgolink

func main() {
	logger.SetLevel(logrus.DebugLevel)
	var err error
	/*dlog, err := dislog.NewDisLogByToken(httpClient, logrus.InfoLevel, logWebhookToken, dislog.InfoLevelAndAbove...)
	if err != nil {
		logger.Errorf("error initializing dislog %s", err)
		return
	}
	defer dlog.Close()

	logger.AddHook(dlog) */
	logger.Infof("starting bot...")

	bot, err = core.NewBot(token,
		core.WithHTTPClient(httpClient),
		core.WithLogger(logger),
		core.WithGatewayConfigOpts(gateway.WithGatewayIntents()),
		core.WithCacheConfigOpts(
			core.WithCacheFlags(core.CacheFlagsDefault|core.CacheFlagVoiceStates),
			core.WithMemberCachePolicy(core.MemberCachePolicyNone),
			core.WithMessageCachePolicy(core.MessageCachePolicyNone),
		),
		core.WithEventListeners(&core.ListenerAdapter{
			OnSlashCommand: onSlashCommand,
			OnButtonClick: onButtonClick,
		}),
	)
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
	}

	dgolink = disgolink.NewDisgolink(bot)
	registerNodes()

	if err = bot.Connect(); err != nil {
		logger.Fatalf("error while connecting to discord: %s", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func registerNodes() {
	dgolink.AddNode(&api.NodeOptions{
		Name:     "kittybot",
		Host:     "lavalink.kittybot.de",
		Port:     "443",
		Password: "8v675n4645804v6839b37c4n6v53897c5",
		Secure:   true,
	})
}
