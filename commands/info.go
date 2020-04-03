package commands

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/tcnksm/go-httpstat"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InfoCommand(session disgord.Session, event *disgord.MessageCreate) {
	if event.Message.Content != "bd!info" {
		return
	} else {
		gatewayPing, _ := session.AvgHeartbeatLatency()
		owner, _ := session.GetUser(event.Ctx, snowflake.ParseSnowflakeString("175610330217447424"))
		self, _ := session.GetCurrentUser(event.Ctx)
		var sb strings.Builder
		sb.WriteString("**REST: " + strconv.Itoa(apiPing("https://discordapp.com/api/v6")) + " ms\n")
		sb.WriteString("WebSocket: " + strconv.Itoa(int(gatewayPing / time.Millisecond)) + " ms\n")
		sb.WriteString("CDN: " + strconv.Itoa(apiPing("https://cdn.discordapp.com")) + "ms**")
		_, err := session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "docbot",
				Description: "**a helper bot written with disgord**",
				Color: 0xd32ce6,
				Fields: []*disgord.EmbedField{
					{
						Name:   "**disgord**",
						Value:  "**version: " + disgord.Version + "**",
						Inline: true,
					},
					{
						Name:   "**author**",
						Value:  "**" + owner.Username + "#" + owner.Discriminator.String() + "**",
						Inline: true,
					},
					{
						Name:   "**ping**",
						Value:  sb.String(),
						Inline: false,
					},
				},
			},
		})

		if err != nil {
			_, _ = session.CreateMessage(event.Ctx, event.Message.ChannelID, &disgord.CreateMessageParams{
				Content: "Failed to send message!\n" + err.Error() ,
			})
		}
	}
}

func apiPing(url string) int {
	if req, err := http.NewRequest("HEAD", url, nil); err != nil {
		log.Fatalln(err)
		return -1
	} else {
		var result httpstat.Result
		ctx := httpstat.WithHTTPStat(req.Context(), &result)
		req = req.WithContext(ctx)
		client := http.DefaultClient

		if res, err := client.Do(req); err != nil {
			log.Fatalln(err)
			return -1
		} else {
			if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
				log.Fatalln(err)
			}
			_ = res.Body.Close()

			return int(result.TCPConnection.Milliseconds())
		}
	}
}
