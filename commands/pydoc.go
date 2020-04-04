package commands

import (
	"github.com/anaskhan96/soup"
	"github.com/andersfylling/disgord"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func PyDocCommand(session disgord.Session, event *disgord.MessageCreate) {
	if args := regexp.MustCompile("\\s+").Split(event.Message.Content, 2); args[0] == "python" {
		standardPyDocs(session, event, args)
	}
}

func standardPyDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/python-3.8.2-docs-html/library", func(path string, info os.FileInfo, err error) error {
		if err == nil && regexp.MustCompile("^((?i)"+args[1]+"\\.html)").MatchString(info.Name()) {
			if c, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				url = "https://docs.python.org/3/library/" + info.Name()
				content = string(c)
				return nil
			}
		} else {
			return err
		}
	}); err != nil {
		_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Content: ReadFailed,
		})
	} else {
		doc := soup.HTMLParse(content)
		title := doc.Find("h1")
		summary := doc.FindAll("p")

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       strings.ReplaceAll(args[1], "_", "\\_") + title.Text(),
				URL:         url,
				Color:       0x3572a5,
				Description: strings.ReplaceAll(summary[0].FullText(), "'", "`"),
			},
		}); err != nil {
			log.Println(err)
		}
	}
}
