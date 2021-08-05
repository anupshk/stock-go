package main

import (
	"fmt"

	"github.com/anupshk/stock/cli"
)

func main() {
	err := cli.Setup()
	defer cli.CancelFunc()
	defer cli.CloseDB()
	if err != nil {
		fmt.Println("Setup Error", err.Error())
		return
	}
	cli.RunApp()
}
