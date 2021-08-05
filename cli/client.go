package cli

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
)

type Client struct {
	Id    int    `pg:",pk"`
	Ident string `survey:"id" pg:",unique"`
	Name  string `survey:"name"`
}

func AddClient() {
	questions := []*survey.Question{
		{
			Name: "ident",
			Prompt: &survey.Input{
				Message: "Client id?",
			},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Client name?",
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}
	client := Client{}
	err := survey.Ask(questions, &client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dbRes, dbErr := InsertClient(&client)
	if dbErr != nil {
		fmt.Println("Error adding client", dbErr)
	} else {
		fmt.Println("Added client", dbRes)
	}
}

func ListClients() {
	fmt.Println("List clients")
	clients := GetAllClients()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})
	for _, v := range clients {
		table.Append([]string{v.Ident, v.Name})
	}
	table.Render()
}

func EditClient() {
	fmt.Println("Edit client")
}

func DeleteClient(client string) {
	fmt.Println("Delete client")
}
