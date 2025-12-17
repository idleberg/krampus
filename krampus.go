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
	Force   bool     `short:"f" help:"Force closing processes by all owners."`
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})
)

func main() {
	ctx := kong.Parse(&CLI)

	// Print help if no arguments are passed
	if isEmpty(CLI.Ports) == 0 && !CLI.Version {
		ctx.PrintUsage(false)
		os.Exit(0)
	}

	switch true {
	case CLI.Version:
		printVersion()

	default:
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
	ports := CLI.Ports

	conns, err := net.Connections("all")
	if err != nil {
		logger.Error("failed to get connections", "error", err)
		return
	}

	for _, port := range ports {
		pid, err := getPIDFromConnections(conns, port)

		if pid == -1 {
			continue
		}

		if err != nil {
			logger.Error(err)
			continue
		}

		err = killProcess(pid, port)

		if err != nil {
			logger.Error(err)
		}
	}
}

func getPIDFromConnections(conns []net.ConnectionStat, port string) (int32, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("invalid port %s", port)
	}

	for _, conn := range conns {
		if conn.Status == "LISTEN" && conn.Laddr.Port == uint32(portInt) {
			return conn.Pid, nil
		}
	}

	logger.Warnf("no process found listening on port %s", port)
	return -1, nil
}

func getPID(port string) (int32, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return 0, err
	}
	return getPIDFromConnections(conns, port)
}

func killProcess(pid int32, port string) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	if !CLI.Force {
		currentUID := os.Getuid()

		uids, err := proc.Uids()
		if err != nil {
			return fmt.Errorf("failed to get process owner: %w", err)
		}

		// Check if current user owns the process (or is root)
		if currentUID != 0 && currentUID != int(uids[0]) {
			return fmt.Errorf("permission denied: process %d is owned by another user. Use --force to override", pid, uids[0], currentUID)
		}
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	logger.Infof("killed process with PID %d, listening on port %s", pid, port)

	return nil
}

func isEmpty(arr []string) int {
	count := 0
	for _, str := range arr {
		if str != "" {
			count++
		}
	}
	return count
}
