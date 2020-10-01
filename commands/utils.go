package commands

import (
	"github.com/andersfylling/disgord"
)

func UtilsCommand(session disgord.Session, event *disgord.MessageCreate) {
	if event.Message.Content != "bd!utils" {
		return
	}
}
