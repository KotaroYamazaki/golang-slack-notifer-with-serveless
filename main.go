package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/pkg/gs"
	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/pkg/slack"
	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/project/config"
)

// 参考：https://github.com/sminamot/sheets-go-example/blob/master/main.go
// https://sminamot-dev.hatenablog.com/entry/2019/12/05/204403

func run() error {

	config := config.GetSecret()

	srv := gs.GetGSConnection(config.Secret)

	memberRow := gs.GetValueRange(srv, config.SheetID, "Sheet1!A2:A")
	nMembers := len(memberRow.Values)
	updateRange := fmt.Sprintf("Sheet1!B2:B%s", strconv.Itoa(nMembers+1))
	facilitatorRow := gs.GetValueRange(srv, config.SheetID, updateRange)

	if (nMembers == 0) || (len(facilitatorRow.Values) == 0) {
		log.Fatal("No data found.")
		return nil
	}

	fIndex := gs.FindFacilitatorIndex(facilitatorRow)
	slack.PostMessage(config.Webhook, fmt.Sprintf("今日の司会は<%s>", memberRow.Values[fIndex][0]))

	// 更新
	_, err := srv.Spreadsheets.Values.Update(config.SheetID, updateRange, gs.Rotate(facilitatorRow, fIndex, nMembers)).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("--- メッセージの通知 完了")
	return nil
}

func main() {
	run()
	//lambda.Start(run)
}
