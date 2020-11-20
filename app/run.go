package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/pkg/gs"
	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/pkg/slack"
	"github.com/KotaroYamazaki/golang-slack-notifer-with-serveless/project/config"
)

func Run() error {

	config := config.GetSecret()

	srv := gs.GetGSConnection(config.Secret)

	memberRow := gs.GetValueRange(srv, config.SheetID, "Sheet1!A2:A")
	nMembers := len(memberRow.Values)
	updateRange := fmt.Sprintf("Sheet1!B2:B%s", strconv.Itoa(nMembers+1))
	facilitatorRow := gs.GetValueRange(srv, config.SheetID, updateRange)

	if (nMembers == 0) || (len(facilitatorRow.Values) == 0) {
		log.Fatal("No data found.")
	}

	fIndex := gs.FindFacilitatorIndex(facilitatorRow)
	slack.PostMessage(config.Webhook, fmt.Sprintf("今日の朝会の司会は%s>\n次回は%sです！", memberRow.Values[fIndex][0], memberRow.Values[(fIndex+1)%nMembers][0]))

	// 更新
	_, err := srv.Spreadsheets.Values.Update(config.SheetID, updateRange, gs.Rotate(facilitatorRow, fIndex, nMembers)).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("--- メッセージの通知 完了")
	return nil
}
