package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/anupshk/stock/util"
	"github.com/olekukonko/tablewriter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Share struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Client       primitive.ObjectID `bson:"client,omitempty" json:"client,omitempty"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Balance      float32            `bson:"balance" json:"balance"`
	PCP          float32            `bson:"prev_close_price" json:"prev_close_price"`
	LTP          float32            `bson:"last_tran_price" json:"last_tran_price"`
	ImportedDate string             `bson:"imported_date" json:"imported_date"`
	ImportedAt   primitive.DateTime `bson:"imported_at" json:"imported_at"`
}

type ShareValue struct {
	Id    time.Time `bson:"_id,omitempty" json:"_id,omitempty"`
	Value float32   `bson:"total" json:"total"`
}

func ListShares(args ...string) {
	ident := args[0]
	date := args[1]
	client, err := GetClient(ident)
	if err != nil {
		fmt.Println("Error: Client not found")
		return
	}
	var shares []Share
	if date == "" {
		shares, err = client.GetLatestShares()
	} else {
		shares, err = client.GetShares(date)
	}
	if err != nil {
		fmt.Println("Error", err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	var pcpVal, ltpVal float32 = 0.0, 0.0
	table.SetHeader([]string{"S.No.", "Date", "Scrip", "Balance", "PCP", "PCP Value", "LTP", "LTP Value"})
	for i, v := range shares {
		pcpV := v.Balance * v.PCP
		ltpV := v.Balance * v.LTP
		pcpVal += pcpV
		ltpVal += ltpV
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			util.GetDisplayDate(v.ImportedDate),
			v.Symbol,
			fmt.Sprintf("%0.2f", v.Balance),
			fmt.Sprintf("%0.2f", v.PCP),
			fmt.Sprintf("%0.2f", pcpV),
			fmt.Sprintf("%0.2f", v.LTP),
			fmt.Sprintf("%0.2f", ltpV),
		})
	}
	table.SetFooter([]string{
		"",
		"",
		"",
		"",
		"Total",
		fmt.Sprintf("%0.2f", pcpVal),
		"",
		fmt.Sprintf("%0.2f", ltpVal),
	})
	table.Render()
}
