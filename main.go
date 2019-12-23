package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New(os.Getenv("CONFIG_SLACK_OAUTH_TOKEN"))

func main() {
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("CONFIG_SLACK_VERIFICATION_TOKEN")}))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err = json.Unmarshal([]byte(body), &r)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				if match := regexp.MustCompile(`([0-9a-zA-Z-]+) を実行していただけませんでしょうか`).FindAllStringSubmatch(ev.Text, -1); match != nil {
					api.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("仕方ありませんね。\n%s を実行します。", match[0][1]), false))
				}
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}
