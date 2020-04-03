package commands

import (
	"github.com/andersfylling/disgord"
)

type CommandHandler interface {
	Handle(session disgord.Session, event *disgord.MessageCreate)
}

var Commands = []interface{}{
	JavaDocCommand.Handle,
}
