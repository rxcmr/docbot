package commands

import (
	"log"
	"regexp"

	"github.com/andersfylling/disgord"
)

func UtilsCommand(session disgord.Session, event *disgord.MessageCreate) {
	if event.Message.Content != "bd!utils" {
		return
	} else {
		args := regexp.MustCompile("\\s+").Split(event.Message.Content, 3)
		if args[1] == "heartbeat" {
			heartbeatAck, _ := session.AvgHeartbeatLatency()
			session.SendMsg(event.Ctx, event.Message.ChannelID, heartbeatAck.String())
		} else if args[1] == "log" {
			log.Println(args[2])
		} else if args[1] == "echo" {
			session.SendMsg(event.Ctx, event.Message.ChannelID, args[2])
		}
	}
}
