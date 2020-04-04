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
		standardJDocs(session, event, args)
	}
}

func standardJDocs(session disgord.Session, event *disgord.MessageCreate, args []string) {
	var content string
	var url string
	if err := filepath.Walk("./resources/docs/api/java.base", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "class-use" {
			return filepath.SkipDir
		} else if err == nil && regexp.MustCompile("^((?i)\\b"+args[1]+"\\.html\\b)").MatchString(info.Name()) &&
			strings.ToLower(args[1]+".html") == strings.ToLower(info.Name()) {
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
	}); err != nil || len(content) == 0 {
		_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Content: ReadFailed,
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
		summary := description.FindStrict("div", "class", "block")

		var tree string
		var definition string
		if workaround, ok := JavaDocWorkarounds[title.Text()]; ok {
			definition = workaround(summary, &tree)
		} else {
			definition = regexp.MustCompile("\n\n").Split(summary.FullText(), -1)[0]
		}

		if inheritanceTree.Pointer == nil && tree == "" {
			tree = "No inheritance tree available. If you see this, it means it wasn't worked around yet."
		} else if inheritanceTree.Pointer != nil {
			tree = inheritanceTree.FullText()
		}

		if _, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       title.Text(),
				URL:         url,
				Color:       0xb07219,
				Description: definition,
				Fields: []*disgord.EmbedField{
					{
						Name:   "Class Signature",
						Value:  additionalInfo.FullText(),
						Inline: true,
					},
					{
						Name:   "Inheritance Tree",
						Value:  tree,
						Inline: true,
					},
				},
			},
		}); err != nil {
			log.Println(err)
		}
	}
}
