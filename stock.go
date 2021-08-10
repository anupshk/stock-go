package main

import (
	"fmt"

	"github.com/anupshk/stock/cli"
)

func main() {
	err := cli.Setup()
	if err != nil {
		fmt.Println("Setup Error", err.Error())
		return
	}
	defer cli.CancelFunc()
	defer cli.CloseDB()
	cli.RunApp()
}
