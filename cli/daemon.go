package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tera-insights/ticrypt-file-copy/config"
	"github.com/tera-insights/ticrypt-file-copy/daemon"
)

func StartDaemon(config *config.Config) error {
	// Start the daemon
	// Try to recover any in progress copies since last time the daemon was running
	err := Recover(config)
	if err != nil {
		return err
	}

	// Create the daemon
	fmt.Println("Starting daemon...")
	fmt.Printf("Listening on port %s\n", config.Server.Port)
	fmt.Printf("Allowed hosts: %v\n", config.Server.AllowedHosts)
	daemon := daemon.NewDaemon(config.Server.Port, config.Server.AllowedHosts)
	// Start the daemon
	err = daemon.Start()
	if err != nil {
		return err
	}

	// Wait for stop signal
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	// Block until a termination signal is received
	<-stopSignal
	fmt.Println("Received termination signal. Shutting down...")

	// Stop the job manager
	daemon.Close()

	return nil
}
