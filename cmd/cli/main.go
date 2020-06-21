package main

import (
	"fmt"
	"os"
	"toad/pkg/get"
	"toad/pkg/post"

	"github.com/urfave/cli"
)

func main() {

	// Create struct that manages flags
	app := cli.NewApp()
	app.Name = "Toad"
	app.Usage = "Load testing tool"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "ap813",
		},
	}

	// Option flags for http get calls
	getFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "url",
			Usage:    "URL of service to hit",
			Value:    "",
			Required: true,
		},
		cli.IntFlag{
			Name:     "vus",
			Usage:    "Number of virtual users",
			Value:    10,
			Required: false,
		},
		cli.IntFlag{
			Name:     "delay",
			Usage:    "Millisecond delay between http calls",
			Value:    100,
			Required: false,
		},
		cli.StringFlag{
			Name:     "headers",
			Usage:    "Headers of call separated by ','",
			Value:    "",
			Required: false,
		},
		cli.IntFlag{
			Name:     "timeout",
			Usage:    "Timeout for single http call",
			Value:    60,
			Required: false,
		},
		cli.IntFlag{
			Name:     "duration",
			Usage:    "Duration of the test in seconds",
			Value:    0,
			Required: true,
		},
		cli.BoolFlag{
			Name:     "debug",
			Usage:    "Debugging option prints http response status code and response body",
			Required: false,
		},
	}

	nonGetFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "url",
			Usage:    "URL of service to hit",
			Value:    "",
			Required: true,
		},
		cli.IntFlag{
			Name:     "vus",
			Usage:    "Number of virtual users, defaults to 10",
			Value:    10,
			Required: false,
		},
		cli.IntFlag{
			Name:     "delay",
			Usage:    "Millisecond delay between http calls",
			Value:    100,
			Required: false,
		},
		cli.StringFlag{
			Name:     "headers",
			Usage:    "Headers of call separated by ','",
			Value:    "",
			Required: false,
		},
		cli.IntFlag{
			Name:     "timeout",
			Usage:    "Timeout for single http call",
			Value:    60,
			Required: false,
		},
		cli.IntFlag{
			Name:     "duration",
			Usage:    "Duration of the test in seconds",
			Value:    0,
			Required: true,
		},
		cli.StringFlag{
			Name:     "body",
			Usage:    "Duration of the test in seconds",
			Value:    "{}",
			Required: true,
		},
		cli.BoolFlag{
			Name:     "debug",
			Usage:    "Debugging option prints http response status code and response body",
			Required: false,
		},
	}

	// Apply commands and their flags
	app.Commands = []cli.Command{
		{
			Name:      "get",
			Usage:     "HTTP Method Get Calls to service",
			UsageText: "Ex. 'toad get -url=localhost:8080/health -d=100 -vus=1 -dur=10'",
			Flags:     getFlags,
			Action: func(c *cli.Context) error {
				return get.HTTPGet(c)
			},
		},
		{
			Name:      "post",
			Usage:     "HTTP Method Post Calls to service",
			UsageText: "Ex. 'toad post -url=localhost:8080/something -d=100 -vus=1 -dur=10 -post=\"{}\"'",
			Flags:     nonGetFlags,
			Action: func(c *cli.Context) error {
				fmt.Println("Post call works")
				return post.HTTPPost(c)
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
