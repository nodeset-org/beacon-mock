package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nodeset-org/beacon-mock/db"
	"github.com/nodeset-org/beacon-mock/server"
	"github.com/urfave/cli/v2"
)

const (
	Version string = "0.1.0"
)

// Run
func main() {
	// Initialise application
	app := cli.NewApp()

	// Set application info
	app.Name = "beacon-mock"
	app.Usage = "Partial mock of a Beacon Chain client, useful for testing applications that use the validator status routes"
	app.Version = Version
	app.Authors = []*cli.Author{
		{
			Name:  "Nodeset",
			Email: "info@nodeset.io",
		},
	}
	app.Copyright = "(C) 2024 NodeSet LLC"

	ipFlag := &cli.StringFlag{
		Name:    "ip",
		Aliases: []string{"i"},
		Usage:   "The IP address to bind the API server to",
		Value:   "127.0.0.1",
	}
	portFlag := &cli.UintFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Usage:   "The port to bind the API server to",
		Value:   48812,
	}
	configFlag := &cli.StringFlag{
		Name:    "config-file",
		Aliases: []string{"c"},
		Usage:   "An optional configuration file to load. If not specified, defaults will be used",
	}

	app.Flags = []cli.Flag{
		ipFlag,
		portFlag,
		configFlag,
	}
	app.Action = func(c *cli.Context) error {
		logger := slog.Default()

		// Load the config file if specified
		config := db.NewDefaultConfig()
		configFile := c.String(configFlag.Name)
		if configFile != "" {
			// Make sure the file exists
			_, err := os.Stat(configFile)
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "Config file [%s] doesn't exist %v", configFile, err)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v", err)
				os.Exit(1)
			}

			// Read the file
			bytes, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v", err)
				os.Exit(1)
			}

			// Unmarshal it
			err = json.Unmarshal(bytes, &config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling config file: %v", err)
				os.Exit(1)
			}
		}

		// Create the server
		var err error
		ip := c.String(ipFlag.Name)
		port := uint16(c.Uint(portFlag.Name))
		server, err := server.NewBeaconMockServer(logger, ip, port, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating server: %v", err)
			os.Exit(1)
		}

		// Start it
		wg := &sync.WaitGroup{}
		err = server.Start(wg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error starting server: %v", err)
			os.Exit(1)
		}
		port = server.GetPort()

		// Handle process closures
		termListener := make(chan os.Signal, 1)
		signal.Notify(termListener, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-termListener
			fmt.Println("Shutting down...")
			err := server.Stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error stopping server: %v", err)
				os.Exit(1)
			}
		}()

		// Run the daemon until closed
		logger.Info(fmt.Sprintf("Started Beacon mock server on %s:%d", ip, port))
		wg.Wait()
		fmt.Println("Server stopped.")
		return nil
	}

	// Run application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
