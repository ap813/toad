package main

import (
	"fmt"
	"os"
	"toad/pkg/get"

	"github.com/urfave/cli"
)

func main() {

	// Create struct that manages flags
	app := cli.NewApp()
	app.Name = "Toad"
	app.Usage = "Load testing tool"

	// Option flags for http get calls
	getFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "url",
			Usage:    "URL of service to hit",
			Value:    "",
			Required: true,
		},
		cli.StringFlag{
			Name:     "headers",
			Usage:    "Headers of call separated by ','",
			Value:    "",
			Required: false,
		},
	}

	// Apply commands and their flags
	app.Commands = []cli.Command{
		{
			Name:      "get",
			Usage:     "HTTP Method Get Calls to services",
			UsageText: "Ex. 'toad get --url=https://google.com",
			Flags:     getFlags,
			Action: func(c *cli.Context) error {
				fmt.Println("Get call works")
				return get.HTTPGet(c)
			},
		},
	}

	// Read the user call and take action accordingly
	err := app.Run(os.Args)

	// Process error
	if err != nil {
		fmt.Println(err)

		// TODO: specific error correlation to status code
		os.Exit(1)
	}

	// For clarity
	os.Exit(0)
}
