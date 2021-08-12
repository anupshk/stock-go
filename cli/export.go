package cli

import (
	"fmt"
	"strconv"

	"github.com/anupshk/stock/util"
)

func (client Client) ExportShares() {
	shares, err := client.GetLatestShares()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	if len(shares) == 0 {
		fmt.Println("No shares found")
		return
	}
	data := [][]string{{"S.N", "Scrip", "Current Balance", "Previous Closing Price", "Value as of Previous Closing Price", "Last Transaction Price (LTP)", "Value as of LTP"}}
	totalPCP, totalLTP := 0.0, 0.0
	for i, row := range shares {
		pcpV := row.Balance * row.PCP
		ltpV := row.Balance * row.LTP
		totalPCP += float64(pcpV)
		totalLTP += float64(ltpV)
		data = append(data, []string{
			strconv.Itoa(i + 1),
			row.Symbol,
			fmt.Sprintf("%0.2f", row.Balance),
			fmt.Sprintf("%0.2f", row.PCP),
			fmt.Sprintf("%0.2f", pcpV),
			fmt.Sprintf("%0.2f", row.LTP),
			fmt.Sprintf("%0.2f", ltpV),
		})
	}
	data = append(data, []string{"Total :", " ", " ", " ", fmt.Sprintf("%0.2f", totalPCP), " ", fmt.Sprintf("%0.2f", totalLTP)})
	err = util.ExportToCsv(data)
	if err != nil {
		fmt.Println("Export error", err)
	} else {
		fmt.Println("Exported successfully")
	}
}
