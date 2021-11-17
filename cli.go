
package main

import (
	"fmt"
	"encoding/json"
	"github.com/urfave/cli/v2"
)

var App = &cli.App{
	Name: "takoyaki",
	Usage: "run and manage virtual private servers",
	Commands: []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "run the takoyaki server",
			Action:  serverAction,
		},
		{
			Name:   "db",
			Aliases: []string{"d"},
			Usage:  "manage the database",
			Subcommands: []*cli.Command{
				{
					Name:   "migrate",
					Usage:  "migrates the database",
					Action: dbMigrateAction,
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
					Action: requestListAction,
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
					Action: requestApproveAction,
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
					Action: requestRejectAction,
				},
			},
		},
	},
}

func serverAction(c *cli.Context) error {

	StartServer()

    return nil
}

func dbMigrateAction(c *cli.Context) error {

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

func requestListAction(c *cli.Context) error {

	db, err := DBConnection()
	if err != nil {
        return cli.Exit("Could not establish connection to database", 1)
	}

	createRequests, err := DBRequestListWithPurpose(db, REQUEST_PURPOSE_VPS_CREATE)
	if err != nil {
        return cli.Exit("Could not fetch create requests from db", 1)
	}

	upgradeRequests, err := DBRequestListWithPurpose(db, REQUEST_PURPOSE_VPS_UPGRADE)
	if err != nil {
        return cli.Exit("Could not fetch upgrade requests from db", 1)
	}

	fmt.Printf("VPS CREATION REQUESTS =-=-=-=-=-=-=\n")
	fmt.Printf("Request ID | Username | RAM | CPU | Disk | OS\n")
	for _, request := range createRequests {

		requestData := VPSCreateRequestData{}
		err := json.Unmarshal([]byte(request.RequestData), &requestData)
		if err != nil {
			return cli.Exit("Error unmarshalling request data", 1)
		}

		fmt.Printf(
			"%d | %s | %d | %d | %d | %s",
			request.ID, request.User.Username, requestData.RAM, requestData.CPU, requestData.Disk, requestData.OS,
		)
	}

	fmt.Printf("VPS UPGRADE REQUESTS =-=-=-=-=-=-=\n")
	fmt.Printf("Request ID | Username | RAM | CPU | Disk\n")
	for _, request := range upgradeRequests {

		requestData := VPSUpgradeRequestData{}
		err := json.Unmarshal([]byte(request.RequestData), &requestData)
		if err != nil {
			return cli.Exit("Error unmarshalling request data", 1)
		}

		fmt.Printf(
			"%d | %s | %d | %d | %d",
			request.ID, request.User.Username, requestData.RAM, requestData.CPU, requestData.Disk,
		)
	}

    return nil
}

func requestApproveAction(c *cli.Context) error {
    return nil
}

func requestRejectAction(c *cli.Context) error {
    return nil
}

