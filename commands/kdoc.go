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

func KDocCommand(session disgord.Session, event *disgord.MessageCreate) {
	if args := regexp.MustCompile("\\s+").Split(event.Message.Content, 2); args[0] == "kotlin" {
		standardKDocs(session, event, args)
	}
}

func standardKDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/kotlin/kotlinlang.org/api/latest/jvm/stdlib", func(path string, info os.FileInfo, err error) error {
		if err == nil && regexp.MustCompile("^((?i)\\b"+args[1]+"\\b)").MatchString(info.Name()) &&
			strings.ToLower(args[1]) == strings.ToLower(info.Name()) && info.IsDir() {
			path += "/index.html"
			if c, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				url = "https://kotlinlang.org/api/latest/jvm/stdlib/" + info.Name() + "/"
				content = string(c)
				return nil
			}
		} else {
			return err
		}
	}); err != nil || len(content) == 0 {
		_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Content: ReadFailed,
		})
	} else {
		doc := soup.HTMLParse(content)
		main := doc.Find("article", "role", "main")
		title := main.Find("h2")
		breadcrumbs := main.Find("div", "class", "api-docs-breadcrumbs").FindAll("a")
		summary := main.Find("p")

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       title.Text(),
				URL:         url,
				Color:       0xf18e33,
				Description: "`" + breadcrumbs[0].Text() + "/" + breadcrumbs[1].Text() + "`\n\n" + summary.FullText(),
			},
		}); err != nil {
			log.Println(err)
		}
	}
}
