package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"github.com/anaskhan96/soup"
	"github.com/andersfylling/disgord"
)

type JavaDocCommand struct{}

func (cmd JavaDocCommand) Handle(session disgord.Session, event *disgord.MessageCreate) {
	if event.Message.Content == "java" {
		var content string
		if err := filepath.Walk("./resources/docs/api/java.base", func(path string, info os.FileInfo, err error) error {
			if err == nil && regexp.MustCompile("((?i)" + event.Message.Content + "\\.html)").MatchString(info.Name()) {
				if c, err := ioutil.ReadFile(info.Name()); err != nil {
					return err
				} else {
					content = string(c)
					return nil
				}
			} else {
				return err
			}
		}); err != nil {
			session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
				Content: READ_FAILED,
			})
		} else {
			doc := soup.HTMLParse(content)
			title := doc.Find("title")
			block := doc.Find("div", "class", "block")
			description := regexp.MustCompile("\\n").Split(block.Text(), -1)

			session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title: title.Text(),
					Description: description[0],
					Thumbnail: &disgord.EmbedThumbnail{
						URL: "https://logos-download.com/wp-content/uploads/2016/10/Java_logo.png",
					},
					Footer: &disgord.EmbedFooter{
						IconURL: "https://www.stickpng.com/assets/images/58480979cef1014c0b5e4901.png",
					},
				},
			})
		}
	}
}
