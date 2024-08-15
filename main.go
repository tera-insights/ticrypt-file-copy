package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/superhawk610/bar"
	"github.com/tera-insights/ticrypt-file-copy/copy"
	"github.com/tera-insights/ticrypt-file-copy/daemon"
	ticrypt "github.com/tera-insights/ticrypt-go"
	"github.com/ttacon/chalk"
	"github.com/urfave/cli/v2"
)

func main() {
	var hostID string = "ticrypt"
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

			progress := make(chan copy.Progress)
			go func() {
				stat := <-progress
				b := bar.NewWithOpts(
					bar.WithDimensions(int(stat.TotalBytes), 20),
					bar.WithFormat(
						fmt.Sprintf(
							" %scopying... %s :percent :bar %s:rate Bytes/s%s :eta",
							chalk.Blue,
							chalk.Reset,
							chalk.Green,
							chalk.Reset,
						),
					),
				)
				for p := range progress {
					b.Update(p.BytesWritten, bar.Context{})
				}
				b.Done()
			}()
			err := copy.NewCopier(source, destination, 4, progress).Copy(copy.Read, copy.Write)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:      "start-daemon",
				Aliases:   []string{"d"},
				Usage:     "Start the daemon",
				UsageText: "start-daemon [host:port]",
				Action: func(c *cli.Context) error {
					// Get the host:port
					host := c.Args().First()
					if host == "" {
						host = "localhost:8080"
					}
					// Create ticrypt client
					tcClient, err := ticrypt.NewClient(&ticrypt.Options{
						Host:  host,
						NoTLS: true,
					})
					if err != nil {
						return err
					}

					err = recover()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return err
					}

					// Create the daemon
					daemon := daemon.NewDaemon(hostID, &tcClient)
					// Start the daemon
					daemon.Start()

					// Wait for stop signal
					stopSignal := make(chan os.Signal, 1)
					signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

					// Block until a termination signal is received
					<-stopSignal
					fmt.Println("Received termination signal. Shutting down...")

					// Stop the job manager
					daemon.Close()

					return nil
				},
			},
			{
				Name:      "recover",
				Aliases:   []string{"r"},
				Usage:     "Recover inturrupted copy",
				UsageText: "recover",
				Action: func(c *cli.Context) error {

					return nil
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

					// Benchmark the copy
					progress := make(chan copy.Progress)
					go func() {
						stat := <-progress
						b := bar.NewWithOpts(
							bar.WithDimensions(int(stat.TotalBytes), 20),
							bar.WithFormat(
								fmt.Sprintf(
									" %s copying...%s :percent :bar %s:rate Bytes/s%s :eta",
									chalk.Blue,
									chalk.Reset,
									chalk.Green,
									chalk.Reset,
								),
							),
						)
						for p := range progress {
							b.Update(p.BytesWritten, bar.Context{})
						}
						b.Done()
					}()
					err := copy.NewCopier(source, destination, 4, progress).Benchmark(copy.Read, copy.Write)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
