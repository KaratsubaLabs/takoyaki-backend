package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"strconv"
)

var App = &cli.App{
	Name:  "takoyaki",
	Usage: "run and manage virtual private servers",
	Commands: []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "run the takoyaki server",
			Action:  serverAction,
		},
		{
			Name:    "db",
			Aliases: []string{"d"},
			Usage:   "manage the database",
			Subcommands: []*cli.Command{
				{
					Name:   "migrate",
					Usage:  "migrates the database",
					Action: dbMigrateAction,
				},
			},
		},
		{
			Name:    "request",
			Aliases: []string{"r"},
			Usage:   "manage user requests",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "list all pending requests",
					Action: requestListAction,
				},
				{
					Name:  "approve",
					Usage: "approve requests",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Usage:   "approve all pending requests",
							Aliases: []string{"a"},
						},
					},
					Action: requestApproveAction,
				},
				{
					Name:  "reject",
					Usage: "reject requests",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Usage:   "reject all pending requests",
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
	fmt.Printf("Request ID | Email | RAM | CPU | Disk | OS\n")
	for _, request := range createRequests {

		requestData := VPSCreateRequestData{}
		err := json.Unmarshal([]byte(request.RequestData), &requestData)
		if err != nil {
			return cli.Exit("Error unmarshalling request data", 1)
		}

		fmt.Printf(
			"%d | %s | %d | %d | %d | %s\n",
			request.ID, request.User.Email, requestData.RAM, requestData.CPU, requestData.Disk, requestData.OS,
		)
	}

	fmt.Printf("VPS UPGRADE REQUESTS =-=-=-=-=-=-=\n")
	fmt.Printf("Request ID | Email | RAM | CPU | Disk\n")
	for _, request := range upgradeRequests {

		requestData := VPSUpgradeRequestData{}
		err := json.Unmarshal([]byte(request.RequestData), &requestData)
		if err != nil {
			return cli.Exit("Error unmarshalling request data", 1)
		}

		fmt.Printf(
			"%d | %s | %d | %d | %d\n",
			request.ID, request.User.Email, requestData.RAM, requestData.CPU, requestData.Disk,
		)
	}

	return nil
}

func requestApproveAction(c *cli.Context) error {

	db, err := DBConnection()
	if err != nil {
		return cli.Exit("Could not establish connection to database", 1)
	}

	if c.Bool("all") {
	}

	if c.NArg() != 1 {
		return cli.Exit("Please pass in only one request ID", 1)
	}

	requestID, err := strconv.ParseUint(c.Args().Get(0), 10, 64)
	if err != nil {
		return cli.Exit("Invalid request ID", 1)
	}

	userRequest, err := DBRequestByID(db, uint(requestID))
	if err != nil {
		return cli.Exit("Error retriving user request", 1)
	}

	// do what needs to be done based on request type
	switch userRequest.RequestPurpose {
	case REQUEST_PURPOSE_VPS_CREATE:

		// parse request data
		requestData := VPSCreateRequestData{}
		err = json.Unmarshal([]byte(userRequest.RequestData), &requestData)
		if err != nil {
			return cli.Exit("Error parsing request data", 1)
		}

		// generate random name for vm
		vmName := RandomString()

		// perform the creation
		err = VPSCreate(vmName, requestData)
		if err != nil {
			return cli.Exit("Failed creating vm", 1)
		}

		// add vps to database
		newVPS := VPS{
			DisplayName:  requestData.DisplayName,
			InternalName: vmName,
			UserID:       requestData.UserID,
			RAM:          requestData.RAM,
			CPU:          requestData.CPU,
			Disk:         requestData.Disk,
			OS:           requestData.OS,
		}
		err = DBVPSCreate(db, newVPS)
		if err != nil {
			return cli.Exit("Failed inserting vm into db", 1)
		}

	default:
		return cli.Exit("Invalid request type", 1)
	}

	// remove the request after it is processed
	err = DBRequestDelete(db, uint(requestID))
	if err != nil {
		return cli.Exit("Error deleting request", 1)
	}

	return nil
}

func requestRejectAction(c *cli.Context) error {

	db, err := DBConnection()
	if err != nil {
		return cli.Exit("Could not establish connection to database", 1)
	}

	if c.Bool("all") {
		err = DBRequestTruncate(db)
		return cli.Exit("Failed to delete all requests", 1)
	}

	if c.NArg() != 1 {
		return cli.Exit("Please pass in only one request ID", 1)
	}

	requestID, err := strconv.ParseUint(c.Args().Get(0), 10, 64)
	if err != nil {
		return cli.Exit("Invalid request ID", 1)
	}

	err = DBRequestDelete(db, uint(requestID))
	if err != nil {
		return cli.Exit("Error deleting request", 1)
	}

	return nil
}
