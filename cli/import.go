package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/anupshk/stock/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Import(ident string, date string) {
	var importedDate time.Time = util.GetCurrentTime()
	var err error
	if date != "" {
		importedDate, err = time.Parse(date, util.DATE_FORMAT)
	}
	if err != nil {
		fmt.Println("Invalid date format", err)
		return
	}
	client, err := GetClient(ident)
	if err != nil {
		fmt.Println("Error: Client not found")
		return
	}
	fmt.Println("Import csv", client)
	files, err := util.GetAllCsv()
	for _, file := range files {
		saveToDB(client, file, importedDate)
	}
}

func saveToDB(client Client, file string, importedDate time.Time) error {
	csvFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer csvFile.Close()
	var rows []Share
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
		return err
	}

	now := primitive.NewDateTimeFromTime(importedDate)
	date := importedDate.Format(util.DATETIME_FORMAT)
	for i, line := range csvLines {
		if i == 0 || i == len(csvLines)-1 {
			continue
		}
		balance, _ := strconv.ParseFloat(line[2], 32)
		pcp, _ := strconv.ParseFloat(line[3], 32)
		ltp, _ := strconv.ParseFloat(line[5], 32)
		share := Share{
			Client:       client.Id,
			Symbol:       line[1],
			Balance:      float32(balance),
			PCP:          float32(pcp),
			LTP:          float32(ltp),
			ImportedDate: date,
			ImportedAt:   now,
		}
		rows = append(rows, share)
	}
	dbRes, dbErr := InsertShares(rows)
	fmt.Println(dbRes, dbErr)
	csvFile.Close()
	e := os.Remove(file)
	if e != nil {
		fmt.Println("Error removing file", e)
	}
	return nil
}
