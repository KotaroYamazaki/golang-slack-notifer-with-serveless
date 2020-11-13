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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	// 参考：https://github.com/sminamot/sheets-go-example/blob/master/main.go
	// https://sminamot-dev.hatenablog.com/entry/2019/12/05/204403
)

const REGION = "ap-northeast-1"

func getSecret() (string, string, string, error) {
	secretName := "dev/slack"
	region := "ap-northeast-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.Fatal(err)
	}

	secretString := aws.StringValue(result.SecretString)

	res := make(map[string]interface{})
	if err := json.Unmarshal([]byte(secretString), &res); err != nil {
		log.Fatal(err)
	}

	return res["SECRET"].(string), res["SHEET_ID"].(string), res["WEBHOOK_URL"].(string), nil
}

func run() error {

	secret, spreadsheetID, IncomingUrl, err := getSecret()
	if err != nil {
		log.Fatal(err)
	}
	srv := getGSConnection(secret)

	memberRow := getValueRange(srv, spreadsheetID, "Sheet1!A2:A")
	n_members := len(memberRow.Values)
	updateRange := fmt.Sprintf("Sheet1!B2:B%s", strconv.Itoa(n_members+1))
	facilitatorRow := getValueRange(srv, spreadsheetID, updateRange)

	if (n_members == 0) || (len(facilitatorRow.Values) == 0) {
		log.Fatal("No data found.")
		return nil
	}

	fIndex := findFacilitatorIndex(facilitatorRow)
	postMessage(IncomingUrl, fmt.Sprintf("今日の司会は<%s>", memberRow.Values[fIndex][0]))

	// 更新
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, updateRange, rotate(facilitatorRow, fIndex, n_members)).ValueInputOption("RAW").Do()
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

func getValueRange(srv *sheets.Service, spreadsheetID string, readRange string) *sheets.ValueRange {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
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
