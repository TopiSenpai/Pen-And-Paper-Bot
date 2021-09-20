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
	"syscall"
)

var sounds = map[discord.Snowflake]map[string]string{
	me: {
		"Wooosh": "QAAAfgIAFFdob29zaCBTb3VuZCBFZmZlY3RzABAyTWlycm9yc0RpYWxvZ3VlAAAAAAAAv2gAC0hEUlZ6d05rVjIwAAEAK2h0dHBzOi8vd3d3LnlvdXR1YmUuY29tL3dhdGNoP3Y9SERSVnp3TmtWMjAAB3lvdXR1YmUAAAAAAAAAAA==",
	},
}

var (
	token       = os.Getenv("token")
	testGuildID = discord.Snowflake(os.Getenv("test_guild_id"))
	me          = discord.Snowflake("170939974227591168")
	httpClient  = http.DefaultClient
	bot         *core.Bot
	dgolink     disgolink.Disgolink
)

func main() {
	log.SetLevel(log.LevelDebug)
	var err error

	log.Infof("starting bot...")

	bot, err = core.NewBot(token,
		core.WithHTTPClient(httpClient),
		core.WithGatewayConfig(gateway.Config{
			GatewayIntents:   discord.GatewayIntentsNonPrivileged | discord.GatewayIntentGuildMembers,
		}),
		core.WithCacheConfigOpts(
			core.WithMemberCachePolicy(core.MemberCachePolicyNone),
			core.WithMessageCachePolicy(core.MessageCachePolicyNone),
		),
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

	if err = bot.Connect(); err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func registerNodes() {
	dgolink.AddNode(&disgolink.NodeOptions{
		Name:     "kittybot",
		Host:     "lavalink.kittybot.de",
		Port:     "443",
		Password: "8v675n4645804v6839b37c4n6v53897c5",
		Secure:   true,
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
		Options: []discord.SlashCommandOption{
			{
				Type:        discord.CommandOptionTypeSubCommand,
				Name:        "add",
				Description: "lets you add new sounds",
				Options: []discord.SlashCommandOption{
					{
						Type:        discord.CommandOptionTypeString,
						Name:        "name",
						Description: "the unique sound name",
						Required:    true,
					},
					{
						Type:        discord.CommandOptionTypeString,
						Name:        "source",
						Description: "the source of the sound(url, ytsearch, etc)",
						Required:    true,
					},
				},
			},
			{
				Type:        discord.CommandOptionTypeSubCommand,
				Name:        "remove",
				Description: "lets you remove existing sounds",
				Options: []discord.SlashCommandOption{
					{
						Type:        discord.CommandOptionTypeString,
						Name:        "name",
						Description: "the unique sound name",
						Required:    true,
					},
				},
			},
			{
				Type:        discord.CommandOptionTypeSubCommand,
				Name:        "list",
				Description: "lets you add new sounds",
			},
		},
		DefaultPermission: true,
	},
}
