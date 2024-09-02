package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/blang/semver"
	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
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

	switch true {
	case CLI.Version:
		printVersion()

	default:
		// Check if all arguments are numbers
		hasValidNumber := false
		for _, port := range CLI.Ports {
			if _, err := strconv.Atoi(port); err == nil {
				hasValidNumber = true
			}
		}

		if !hasValidNumber {
			ctx.PrintUsage(false)
			os.Exit(1)
		}

		killPorts()
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

func killPorts() {
	ports := os.Args[1:]

	for _, port := range ports {
		pid, err := getPID(port)

		if pid == -1 {
			continue
		}

		if err != nil {
			logger.Error(err)
			continue
		}

		err = killProcess(pid)

		if err != nil {
			logger.Error(err)
		}
	}
}

func getPID(port string) (int32, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return 0, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("invalid port \"%s\"", port)
	}

	for _, conn := range conns {
		if conn.Status == "LISTEN" && conn.Laddr.Port == uint32(portInt) {
			return conn.Pid, nil
		}
	}

	logger.Warnf("no process found listening on port %s", port)
	return -1, nil
}

func killProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	logger.Infof("killed process with PID %d", pid)

	return nil
}
