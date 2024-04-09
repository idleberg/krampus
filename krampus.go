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
	Ports []string `arg:"" default:""`
	Version bool `short:"v" help:"Show version."`
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})
)

func main() {
	kong.Parse(&CLI)

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
		return 0, fmt.Errorf("invalid port %s", port)
	}

	for _, conn := range conns {
		if conn.Status == "LISTEN" && conn.Laddr.Port == uint32(portInt) {
			return conn.Pid, nil
		}
	}

	logger.Warn(fmt.Sprintf("no process found listening on port %s", port))
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

	logger.Info(fmt.Sprintf("killed process with PID %d", pid))

	return nil
}
