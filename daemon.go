package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tera-insights/ticrypt-file-copy/daemon"
)

func startDaemon(config *config) error {
	// Start the daemon
	// Try to recover any in progress copies since last time the daemon was running
	err := recover(config)
	if err != nil {
		return err
	}

	// Create the daemon
	daemon := daemon.NewDaemon(config.Server.Port, config.Server.AllowedHosts)
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
}
