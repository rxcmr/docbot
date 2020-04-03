package commands

import (
	"os"

	"github.com/andersfylling/disgord"
)

func ShutdownCommand(session disgord.Session, event *disgord.MessageCreate) {
	owner, _ := session.GetUser(event.Ctx, disgord.ParseSnowflakeString("175610330217447424"))
	process, _ := os.FindProcess(os.Getpid())

	if event.Message.Author != owner {
		return
	} else {
		_ = process.Signal(os.Interrupt)
	}
}
