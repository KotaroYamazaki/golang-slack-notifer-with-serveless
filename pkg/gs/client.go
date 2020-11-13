package gs

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func GetGSConnection(secret string) *sheets.Service {

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

func GetValueRange(srv *sheets.Service, sheetID string, readRange string) *sheets.ValueRange {
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func FindFacilitatorIndex(s *sheets.ValueRange) int {
	index := 0
	for _, row := range s.Values {
		if row[0] == "this" {
			return index
		}
		index += 1
	}
	return -1
}

func Rotate(s *sheets.ValueRange, fIndex int, sum int) *sheets.ValueRange {
	s.Values[fIndex], s.Values[(fIndex+1)%sum] = s.Values[(fIndex+1)%sum], s.Values[fIndex]
	return s
}
