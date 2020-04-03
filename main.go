package main

import (
	"context"
	"docbot/commands"
	"log"
	"os"
	"strconv"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./resources/.env"); err != nil {
		log.Println("Failed to load .env file, using global env vars...")
	}

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
		ProjectName: "docbot",
		ShardConfig: disgord.ShardConfig{
			ShardIDs: []uint{0, 1},
			ShardCount: 2,
		},
		Presence: &disgord.UpdateStatusPayload{
			AFK: true,
			Status: disgord.StatusDnd,
			Game: &disgord.Activity{
				Name: "with documentation",
				Type: 1,
			},
		},
	})
	defer client.StayConnectedUntilInterrupted(context.Background())
	logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix("doc ")
	logger := client.Logger()
	myself, _ := client.GetCurrentUser(context.Background())

	client.Ready(func() {
		inviteURL, _ := client.InviteURL(context.Background())
		logger.Info("disgord instance ready.")
		logger.Info("logged in as: ", myself.Username, "#", myself.Discriminator.String())
		logger.Info("pid: ", os.Getpid())
		logger.Info("OAuth2 Invite URL: ", inviteURL)
	})

	client.GuildsReady(func() {
		if guildCount := len(client.GetConnectedGuilds()); guildCount == 1 {
			logger.Info(strconv.Itoa(guildCount) + " guild available.")
		} else {
			logger.Info(strconv.Itoa(guildCount) + " guilds available.")
		}
	})

	events := []interface{}{
		filter.NotByBot,
		filter.HasPrefix,
		logFilter.LogMsg,
		std.CopyMsgEvt,
		filter.StripPrefix,
	}

	for _, command := range commands.Commands {
		events = append(events, command)
	}

	client.On(disgord.EvtMessageCreate, events...)
	client.On(disgord.EvtGuildMemberAdd, func(session disgord.Session, event *disgord.GuildMemberAdd) {
		self, _ := client.GetCurrentUser(event.Ctx)
		if event.Member.User == self {
			channels, _ := session.GetGuildChannels(event.Ctx, event.Member.GuildID)
			_, _ = session.CreateMessage(event.Ctx, channels[0].ID, &disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title: "docbot",
					Description: "**documentation is love, documentation is life**",
					Thumbnail: &disgord.EmbedThumbnail{
						URL: self.Avatar,
					},
					Fields: []*disgord.EmbedField{{
						Name: "**start by typing doc help**",
						Value: "**made with <3, the go programming language, and andersfylling/disgord**",
						Inline: true,
					}},
					Footer: &disgord.EmbedFooter{
						IconURL: self.Avatar,
						Text: self.Username + "#" + self.Discriminator.String(),
					},
				},
			})
		}
	})
}
