package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anupshk/stock/util"
	"github.com/olekukonko/tablewriter"
	"github.com/wcharczuk/go-chart/v2"
)

func StockValueSummaryGraph(c string, scrip string) error {
	client, err := GetClient(c)
	if err != nil {
		fmt.Println("Error: Client not found")
		return err
	}
	values, err := client.GetStockValueSummary(scrip)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	if len(values) < 1 {
		fmt.Println("No stock")
		return nil
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Value", "Remark"})
	prevVal := 0.0
	valuesArr := []float64{}
	datesArr := []time.Time{}
	for _, v := range values {
		valuesArr = append(valuesArr, float64(v.Value))
		datesArr = append(datesArr, v.Id)
		var diff float32
		if prevVal > 0 {
			diff = v.Value - float32(prevVal)
		}
		diffColor := tablewriter.FgHiGreenColor
		if diff < 0 {
			diffColor = tablewriter.FgRedColor
		}
		table.Rich([]string{
			v.Id.Format(util.DATE_FORMAT),
			fmt.Sprintf("%0.2f", v.Value),
			fmt.Sprintf("%0.2f", diff),
		}, []tablewriter.Colors{{}, {}, {tablewriter.Bold, diffColor}})
		prevVal = float64(v.Value)
	}
	table.Render()

	if len(values) < 2 {
		fmt.Println("Not enough data to generate graph")
		return nil
	}

	min, max := util.MinMax(valuesArr)

	priceSeries := chart.TimeSeries{
		Name: "SP",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: datesArr,
		YValues: valuesArr,
	}

	title := "Stock value summary for " + client.Name
	outfile := fmt.Sprintf("%v-%d-summary.png", client.Ident, time.Now().Unix())
	if scrip != "" {
		title = "(" + strings.ToUpper(scrip) + ") value summary for " + client.Name
		outfile = fmt.Sprintf("%v-%v-%d-summary.png", client.Ident, strings.ToUpper(scrip), time.Now().Unix())
	}

	graph := chart.Chart{
		Title: title,
		XAxis: chart.XAxis{
			Name:         "Date",
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Name: "Value",
			Range: &chart.ContinuousRange{
				Max: max,
				Min: min,
			},
		},
		Series: []chart.Series{
			priceSeries,
		},
	}
	f, _ := os.Create("graphs/" + outfile)
	defer f.Close()
	graph.Render(chart.PNG, f)
	return nil
}
