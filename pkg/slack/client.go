package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Slack struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Channel   string `json:"channel"`
}

func PostMessage(IncomingUrl string, msg string) {
	resp, err := http.PostForm(
		IncomingUrl,
		url.Values{"payload": {string(GetPayload(msg))}},
	)
	if err != nil {
		log.Fatal("HTTPリクエストに失敗しました。, err:" + fmt.Sprint(err))
	}

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Http Status:%s, result: %s\n", resp.Status, contents)
}

func GetPayload(msg string) []byte {
	params := Slack{
		Text:      fmt.Sprintf("%s", msg),
		Username:  "From golang to slack hello",
		IconEmoji: ":gopher:",
		IconURL:   "",
		Channel:   "",
	}
	jsonparams, _ := json.Marshal(params)

	return jsonparams
}
