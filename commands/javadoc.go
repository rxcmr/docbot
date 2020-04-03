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

func JavaDocCommand(session disgord.Session, event *disgord.MessageCreate) {
	if args := regexp.MustCompile("\\s+").Split(event.Message.Content, 2); args[0] == "java" {
		standardDocs(session, event, args)
	}
}

func standardDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/docs/api/java.base", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "class-use" {
			return filepath.SkipDir
		} else if err == nil && regexp.MustCompile("^((?i)" + args[1] + "\\.html)").MatchString(info.Name()) {
			if c, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				url = "https://docs.oracle.com/en/java/javase/14/" + strings.ReplaceAll(path, "resources/", "")
				content = string(c)
				return nil
			}
		} else {
			return err
		}
	}); err != nil {
		_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Content: READ_FAILED,
		})
	} else {

		doc := soup.HTMLParse(content)
		main := doc.Find("main", "role", "main")
		header := main.Find("div", "class", "header")
		title := header.Find("h1", "class", "title")
		contentContainer := main.Find("div", "class", "contentContainer")
		inheritanceTree := contentContainer.Find("div", "class", "inheritance")
		description := contentContainer.Find("section", "class", "description")
		additionalInfo := contentContainer.Find("pre")
		summary := description.Find("div", "class", "block")

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:  title.Text(),
				URL: url,
				Color: 0xd32ce6,
				Description: regexp.MustCompile("\n\n").Split(summary.FullText(), -1)[0],
				Fields: []*disgord.EmbedField{
					{
						Name: "Class Signature",
						Value: additionalInfo.FullText(),
						Inline: true,
					},
					{
						Name: "Inheritance Tree",
						Value: inheritanceTree.FullText(),
						Inline: true,
					},
				},
			},
		}); err != nil {
			log.Println(err)
		}
	}
}