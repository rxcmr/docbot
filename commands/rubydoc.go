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

func RubyDocCommand(session disgord.Session, event *disgord.MessageCreate) {
	if args := regexp.MustCompile("\\s+").Split(event.Message.Content, 2); args[0] == "ruby" {
		standardCoreRubyDocs(session, event, args)
	}
}

func standardCoreRubyDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/ruby/ruby_2_7_1_core", func(path string, info os.FileInfo, err error) error {
		if err == nil && regexp.MustCompile("^((?i)\\b"+args[1]+".html\\b)").MatchString(info.Name()) &&
			strings.ToLower(args[1]+".html") == strings.ToLower(info.Name()) {
			if c, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				url = "https://ruby-doc.org/core-2.7.1/" + info.Name()
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
		documentation := doc.Find("div", "id", "documentation")
		title := documentation.Find("h1")
		description := documentation.Find("div", "id", "description")

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       title.Text(),
				URL:         url,
				Color:       0x701516,
				Description: description.Find("p").FullText(),
			},
		}); err != nil {
			log.Println(err)
		}
	}
}
