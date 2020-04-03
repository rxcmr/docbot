package commands

import (
	"github.com/anaskhan96/soup"
	"github.com/andersfylling/disgord"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func GoDocCommand(session disgord.Session, event *disgord.MessageCreate) {
	if args := regexp.MustCompile("\\s+").Split(event.Message.Content, 2); args[0] == "go" {
		standardGoDocs(session, event, args)
	}
}

func standardGoDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/pkg", func(path string, info os.FileInfo, err error) error {
		if err == nil && regexp.MustCompile("^((?i)"+args[1]+")").MatchString(info.Name()) {
			path += "/index.html"
			if c, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				log.Println(path)
				url = "https://godoc.org/" + info.Name()
				content = string(c)
				return nil
			}
		} else {
			return err
		}
	}); err != nil {
		log.Println(err)
		_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Content: ReadFailed,
		})
	} else {
		doc := soup.HTMLParse(content)
		pkgOverview := doc.Find("div", "id", "pkg-overview")
		expanded := pkgOverview.Find("div", "class", "expanded")
		title := doc.Find("h1")
		shortNav := doc.Find("div", "id", "short-nav")
		importStmt := "```go\n" + shortNav.Find("dl").Find("dd").Find("code").Text() + "\n```"
		description := expanded.FindAll("p")[0]

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       title.Text(),
				URL:         url,
				Color:       0x00ADD8,
				Description: importStmt + description.Text(),
			},
		}); err != nil {
			log.Println(err)
		}
	}
}
