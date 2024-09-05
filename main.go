package main

import (
	"fmt"
	"log"
	"os"

	ticli "github.com/tera-insights/ticrypt-file-copy/cli"
	"github.com/tera-insights/ticrypt-file-copy/config"
	"github.com/urfave/cli/v2"
)

func main() {
	// Read the config
	config := config.FetchConfig()
	// Create the CLI app
	app := &cli.App{
		Name:      "ticrypt-file-copy",
		Usage:     "Hight performance tool to copy files",
		UsageText: "ticp [source] [destination]",
		Action: func(c *cli.Context) error {
			// Get the source and destination
			source := c.Args().First()
			if source == "" {
				fmt.Println("Source file is required")
				return nil
			}
			destination := c.Args().Get(1)
			if destination == "" {
				fmt.Println("Destination file is required")
				return nil
			}
			// Copy the file
			return ticli.Ticp(source, destination, config)
		},

		Commands: []*cli.Command{
			{
				Name:      "start-daemon",
				Aliases:   []string{"d"},
				Usage:     "Start the daemon",
				UsageText: "start-daemon [port]",
				Action: func(c *cli.Context) error {
					// Get the port
					port := c.Args().First()
					if port != "" {
						config.Server.Port = port
					}
					return ticli.StartDaemon(config)
				},
			},
			{
				Name:      "recover",
				Aliases:   []string{"r"},
				Usage:     "Recover inturrupted copy",
				UsageText: "recover",
				Action: func(c *cli.Context) error {
					return ticli.Recover(config)
				},
			},
			{
				Name:      "benchmark",
				Aliases:   []string{"b"},
				Usage:     "Run the benchmark",
				UsageText: "benchmark",
				Action: func(c *cli.Context) error {
					// Get the source and destination
					source := c.Args().First()
					if source == "" {
						source = "source"
					}
					destination := c.Args().Get(1)
					if destination == "" {
						destination = "destination"
					}
					return ticli.Benchmark(source, destination, config)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
