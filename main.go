package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgolink"
	"github.com/DisgoOrg/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	token             = os.Getenv("token")
	lavalinkHost      = os.Getenv("lavalink_host")
	lavalinkPort      = os.Getenv("lavalink_port")
	lavalinkPassword  = os.Getenv("lavalink_password")
	lavalinkSecure, _ = strconv.ParseBool(os.Getenv("lavalink_secure"))
	testGuildID       = discord.Snowflake(os.Getenv("test_guild_id"))
	httpClient        = http.DefaultClient
	bot               *core.Bot
	dgolink           disgolink.Disgolink
)

func main() {
	log.SetLevel(log.LevelDebug)
	var err error

	log.Infof("starting bot...")

	loadSounds()

	bot, err = core.NewBot(token,
		core.WithHTTPClient(httpClient),
		core.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentsNonPrivileged)),
		core.WithCacheConfigOpts(core.WithMemberCachePolicy(core.MemberCachePolicyVoice)),
		core.WithEventListeners(&core.ListenerAdapter{
			OnSlashCommand: onSlashCommand,
			OnButtonClick:  onButtonClick,
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
	}

	defer bot.Close()

	if testGuildID == "" {
		_, _ = bot.SetCommands(commands)
	} else {
		_, _ = bot.SetGuildCommands(testGuildID, commands)
	}

	dgolink = disgolink.NewDisgolink(bot)
	registerNodes()

	if err = bot.ConnectGateway(); err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
	saveSounds()
}

func registerNodes() {
	dgolink.AddNode(&disgolink.NodeOptions{
		Name:     "node1",
		Host:     lavalinkHost,
		Port:     lavalinkPort,
		Password: lavalinkPassword,
		Secure:   lavalinkSecure,
	})
}

var commands = []discord.ApplicationCommandCreate{
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "soundboard",
		Description:       "shows your soundboard",
		DefaultPermission: true,
	},
	{
		Type:        discord.ApplicationCommandTypeSlash,
		Name:        "sounds",
		Description: "lets you add/remove/list sounds",
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeSubCommand,
				Name:        "add",
				Description: "lets you add new sounds",
				Options: []discord.ApplicationCommandOption{
					{
						Type:        discord.ApplicationCommandOptionTypeString,
						Name:        "name",
						Description: "the unique sound name",
						Required:    true,
					},
					{
						Type:        discord.ApplicationCommandOptionTypeString,
						Name:        "source",
						Description: "the source of the sound(url, yt, sc, etc)",
						Required:    true,
					},
				},
			},
			{
				Type:        discord.ApplicationCommandOptionTypeSubCommand,
				Name:        "remove",
				Description: "lets you remove existing sounds",
				Options: []discord.ApplicationCommandOption{
					{
						Type:        discord.ApplicationCommandOptionTypeString,
						Name:        "name",
						Description: "the unique sound name",
						Required:    true,
					},
				},
			},
		},
		DefaultPermission: true,
	},
}
