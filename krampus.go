package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/net"
)

type ConnectionLister interface {
	Connections() ([]net.ConnectionStat, error)
}

type ProcessKiller interface {
	Kill(pid int32) error
	Uids(pid int32) ([]int32, error)
}

type Krampus struct {
	Ports  []string
	Force  bool
	UID    int
	Lister ConnectionLister
	Killer ProcessKiller
	Logger *log.Logger
}

func (k *Krampus) Run() {
	conns, err := k.Lister.Connections()
	if err != nil {
		k.Logger.Error("failed to get connections", "error", err)
		return
	}

	for _, port := range k.Ports {
		pid, err := getPIDFromConnections(conns, port)

		if pid == -1 {
			k.Logger.Warnf("no process found listening on port %s", port)
			continue
		}

		if err != nil {
			k.Logger.Error(err)
			continue
		}

		err = k.killProcess(pid, port)
		if err != nil {
			k.Logger.Error(err)
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

	return -1, nil
}

func canKill(currentUID, processUID int, force bool) bool {
	return force || currentUID == 0 || currentUID == processUID
}

func (k *Krampus) killProcess(pid int32, port string) error {
	if !k.Force {
		uids, err := k.Killer.Uids(pid)
		if err != nil {
			return fmt.Errorf("failed to get process owner: %w", err)
		}

		if !canKill(k.UID, int(uids[0]), false) {
			return fmt.Errorf("permission denied: process %d on port %s is owned by UID %d, current UID is %d. Use --force to override", pid, port, uids[0], k.UID)
		}
	}

	err := k.Killer.Kill(pid)
	if err != nil {
		return err
	}

	k.Logger.Infof("killed process with PID %d, listening on port %s", pid, port)
	return nil
}
