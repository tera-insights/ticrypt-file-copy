package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "ticrypt-file-copy",
		Usage:     "Hight performance tool to copy files",
		UsageText: "ticp [source] [destination]",
		Action: func(*cli.Context) error {
			copier := NewCopier("source", "destination", 4)
			start := time.Now()
			err := copier.Copy()
			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Printf("Time taken %v /GB \n", elapsed)
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
				UsageText: "start-daemon",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting the daemon")
					return nil
				},
			},
			{
				Name:      "benchmark",
				Aliases:   []string{"b"},
				Usage:     "Run the benchmark",
				UsageText: "benchmark",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("dd", []string{"if=/dev/urandom", "of=source", "bs=64M", "count=16", "iflag=fullblock"}...)
					start := time.Now()
					err := cmd.Run()
					t := time.Now()
					elapsed := t.Sub(start)
					fmt.Printf("Time taken %v /GB \n", elapsed)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return nil
					}
					fmt.Println("File created")
					defer func() {
						cmd = exec.Command("rm", "source")
						err = cmd.Run()
						if err != nil {
							fmt.Printf("Error: %v\n", err)
						}
						cmd = exec.Command("rm", "destination")
						err = cmd.Run()
						if err != nil {
							fmt.Printf("Error: %v\n", err)
						}
					}()

					copier := NewCopier("source", "destination", 4096)
					fmt.Println("Starting the benchmark")
					start = time.Now()
					err = copier.Copy()
					t = time.Now()
					elapsed = t.Sub(start)
					fmt.Printf("Time taken %v /GB \n", elapsed)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}

					fmt.Println("Benchmark rsync")
					cmd = exec.Command("rsync", []string{"source", "destination"}...)
					start = time.Now()
					err = cmd.Run()
					t = time.Now()
					elapsed = t.Sub(start)
					fmt.Printf("Time taken %v /GB \n", elapsed)
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
