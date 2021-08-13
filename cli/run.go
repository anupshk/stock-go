package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func RunApp() error {
	app := &cli.App{
		Name:  "Stock",
		Usage: "View stock values",
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Value:   "table",
			Usage:   "View type",
			Aliases: []string{"t"},
		},
	}
	app.Action = func(ctx *cli.Context) error {
		if ctx.String("type") == "graph" {
			fmt.Println("Not implemented")
		} else {
			ShowTable()
		}
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:  "migrate",
			Usage: "Migrate database",
			Action: func(ctx *cli.Context) error {
				err := Migrate()
				if err != nil {
					fmt.Println("Error migrating database")
				} else {
					fmt.Println("Database migrated to latest version")
				}
				return nil
			},
		},
		{
			Name:  "graph",
			Usage: "Graph View",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "client",
					Usage:    "Client ID",
					Aliases:  []string{"c"},
					Required: true,
				},
				&cli.StringFlag{
					Name:    "scrip",
					Aliases: []string{"s"},
				},
			},
			Action: func(ctx *cli.Context) error {
				client := ctx.String("client")
				scrip := ctx.String("scrip")
				StockValueSummaryGraph(client, scrip)
				return nil
			},
		},
		{
			Name:    "import",
			Usage:   "Import stock csv",
			Aliases: []string{"i"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "client",
					Usage:    "Client ID",
					Aliases:  []string{"c"},
					Required: true,
				},
				&cli.StringFlag{
					Name:     "date",
					Usage:    "Date",
					Aliases:  []string{"d"},
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				client := ctx.String("client")
				if client == "" {
					fmt.Println("Client ID required")
					return nil
				}
				date := ctx.String("date")
				Import(client, date)
				return nil
			},
		},
		{
			Name:    "export",
			Usage:   "Export stock csv",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "client",
					Usage:    "Client ID",
					Aliases:  []string{"c"},
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				c := ctx.String("client")
				if c == "" {
					fmt.Println("Client ID required")
					return nil
				}
				client, err := GetClient(c)
				if err != nil {
					fmt.Println("Client not found")
					return nil
				}
				client.ExportShares()
				return nil
			},
		},
		{
			Name:    "shares",
			Usage:   "List shares",
			Aliases: []string{"s"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "client",
					Usage:    "Client ID",
					Aliases:  []string{"c"},
					Required: true,
				},
				&cli.StringFlag{
					Name:     "date",
					Usage:    "Imported Date (eg.: 2021-01-25)",
					Aliases:  []string{"d"},
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				client := ctx.String("client")
				if client == "" {
					fmt.Println("Client ID required")
					return nil
				}
				date := ctx.String("date")
				ListShares(client, date)
				return nil
			},
		},
		{
			Name:    "client",
			Usage:   "Client options",
			Aliases: []string{"c"},
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Usage:   "Add client",
					Aliases: []string{"a"},
					Action: func(ctx *cli.Context) error {
						AddClient()
						return nil
					},
				},
				{
					Name:    "list",
					Usage:   "List clients",
					Aliases: []string{"l"},
					Action: func(ctx *cli.Context) error {
						ListClients()
						return nil
					},
				},
				{
					Name:    "edit",
					Usage:   "Edit client",
					Aliases: []string{"e"},
					Action: func(ctx *cli.Context) error {
						EditClient()
						return nil
					},
				},
				{
					Name:    "delete",
					Usage:   "Delete client",
					Aliases: []string{"d", "del"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "client",
							Usage:    "Client ID",
							Aliases:  []string{"c"},
							Required: true,
						},
					},
					Action: func(ctx *cli.Context) error {
						client := ctx.String("client")
						if client == "" {
							fmt.Println("Client ID required")
							return nil
						}
						DeleteClient(client)
						return nil
					},
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	return app.Run(os.Args)
}
