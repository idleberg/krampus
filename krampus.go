package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/blang/semver"
	"github.com/charmbracelet/log"
)

var Version string

var CLI struct {
	Ports   []string `arg:"" default:""`
	Version bool     `short:"v" help:"Show version."`
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})
)

func main() {
	ctx := kong.Parse(&CLI)

	// Print help if no arguments are passed
	if len(CLI.Ports) == 0 && !CLI.Version {
		ctx.PrintUsage(false)
		os.Exit(0)
	}

	switch true {
	case CLI.Version:
		printVersion()

	default:
		// Check if all arguments are numbers
		hasValidNumber := false
		for _, port := range CLI.Ports {
			if _, err := strconv.Atoi(port); err == nil {
				hasValidNumber = true
			} else {
				logger.Warnf("Ignoring invalid port: %s. Error: %v", port, err)
			}
		}

		if !hasValidNumber {
			logger.Error("No valid ports provided, printing `krampus --help`\n")
			ctx.PrintUsage(false)
			os.Exit(1)
		}

		// Filter out non-number arguments
		validPorts := []string{}
		for _, port := range CLI.Ports {
			if _, err := strconv.Atoi(port); err == nil {
				validPorts = append(validPorts, port)
			}
		}

		// Proceed with valid ports
		for _, port := range validPorts {
			fmt.Printf("Processing port: %s\n", port)
		}
	}
}

func printVersion() {
	ver, err := semver.Parse(Version)

	var outputVersion string

	if err == nil {
		outputVersion = "v" + Version
	} else {
		outputVersion = ver.String()
	}

	fmt.Println(outputVersion)
}
