
package main

import (
	"os"
	"github.com/urfave/cli/v2"
)

func main() {

    app := &cli.App{
        Name: "takoyaki",
        Usage: "run and manage virtual private servers",
        Commands: []*cli.Command{
            {
                Name:    "server",
                Aliases: []string{"s"},
                Usage:   "run the takoyaki server",
                Action:  serverHandler,
            },
            {
                Name:   "db",
                Aliases: []string{"d"},
                Usage:  "manage the database",
				Subcommands: []*cli.Command{
					{
						Name:   "migrate",
						Usage:  "migrates the database",
						Action: dbMigrateHandler,
					},
                },
            },
            {
                Name:   "request",
                Aliases: []string{"r"},
                Usage:  "manage user requests",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all pending requests",
						Action: requestListHandler,
					},
					{
						Name:   "approve",
						Usage:  "approve requests",
						Flags:  []cli.Flag{
							&cli.BoolFlag{
								Name: "all",
								Usage: "approve all pending requests",
                                Aliases: []string{"a"},
							},
						},
						Action: requestApproveHandler,
					},
					{
						Name:   "reject",
						Usage:  "reject requests",
						Flags:  []cli.Flag{
							&cli.BoolFlag{
								Name: "all",
								Usage: "reject all pending requests",
                                Aliases: []string{"a"},
							},
						},
						Action: requestRejectHandler,
					},
				},
            },
        },
    }

	app.Run(os.Args)


}

func serverHandler(c *cli.Context) error {

	StartServer()

    return nil
}

func dbMigrateHandler(c *cli.Context) error {

	db, err := DBConnection()
	if err != nil {
        return cli.Exit("Could not establish connection to database", 1)
	}

    err = DBMigrate(db)
	if err != nil {
        return cli.Exit("Migration failed", 1)
	}

    return nil
}

func requestListHandler(c *cli.Context) error {
    return nil
}

func requestApproveHandler(c *cli.Context) error {
    return nil
}

func requestRejectHandler(c *cli.Context) error {
    return nil
}

