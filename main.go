package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	// 参考：https://github.com/sminamot/sheets-go-example/blob/master/main.go
	// https://sminamot-dev.hatenablog.com/entry/2019/12/05/204403
)

const REGION = "ap-northeast-1"

func run() error {

	config := getSecret()
	srv := getGSConnection(config.Secret)

	memberRow := getValueRange(srv, config.SheetID, "Sheet1!A2:A")
	n_members := len(memberRow.Values)
	updateRange := fmt.Sprintf("Sheet1!B2:B%s", strconv.Itoa(n_members+1))
	facilitatorRow := getValueRange(srv, config.SheetID, updateRange)

	if (n_members == 0) || (len(facilitatorRow.Values) == 0) {
		log.Fatal("No data found.")
		return nil
	}

	fIndex := findFacilitatorIndex(facilitatorRow)
	postMessage(config.Webhook, fmt.Sprintf("今日の司会は<%s>", memberRow.Values[fIndex][0]))

	// 更新
	_, err := srv.Spreadsheets.Values.Update(config.SheetID, updateRange, rotate(facilitatorRow, fIndex, n_members)).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("--- メッセージの通知 完了")
	return nil
}

func postMessage(IncomingUrl string, msg string) {
	resp, err := http.PostForm(
		IncomingUrl,
		url.Values{"payload": {string(getPayload(msg))}},
	)
	if err != nil {
		log.Fatal("HTTPリクエストに失敗しました。, err:" + fmt.Sprint(err))
	}

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Http Status:%s, result: %s\n", resp.Status, contents)
}

func getPayload(msg string) []byte {
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

type Slack struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Channel   string `json:"channel"`
}

func getGSConnection(secret string) *sheets.Service {

	conf, err := google.JWTConfigFromJSON([]byte(secret), sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background())
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

func getValueRange(srv *sheets.Service, sheetID string, readRange string) *sheets.ValueRange {
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func findFacilitatorIndex(s *sheets.ValueRange) int {
	index := 0
	for _, row := range s.Values {
		if row[0] == "this" {
			return index
		}
		index += 1
	}
	return -1
}

func rotate(s *sheets.ValueRange, fIndex int, sum int) *sheets.ValueRange {
	s.Values[fIndex], s.Values[(fIndex+1)%sum] = s.Values[(fIndex+1)%sum], s.Values[fIndex]
	return s
}

func main() {
	run()
	//lambda.Start(run)
}
