package main

import (
	"fmt"
	"log"
	"strconv"
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

func main() {
	run()
	//lambda.Start(run)
}
