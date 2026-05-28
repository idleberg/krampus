package main

import (
	"fmt"
	"os"

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

func main() {
	ctx := kong.Parse(&CLI)

	if !CLI.Version && !hasNonEmptyArgs(CLI.Ports) {
		ctx.PrintUsage(false)
		os.Exit(0)
	}

	if CLI.Version {
		printVersion()
		return
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})

	k := &Krampus{
		Ports:  CLI.Ports,
		Force:  CLI.Force,
		UID:    os.Getuid(),
		Lister: &realConnectionLister{},
		Killer: &realProcessKiller{},
		Logger: logger,
	}

	k.Run()
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

func hasNonEmptyArgs(args []string) bool {
	for _, s := range args {
		if s != "" {
			return true
		}
	}
	return false
}

type realConnectionLister struct{}

func (r *realConnectionLister) Connections() ([]net.ConnectionStat, error) {
	return net.Connections("all")
}

type realProcessKiller struct{}

func (r *realProcessKiller) Kill(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func (r *realProcessKiller) Uids(pid int32) ([]int32, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	return proc.Uids()
}
