package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/net"
)

func mockConnection(port uint32, pid int32, status string) net.ConnectionStat {
	return net.ConnectionStat{
		Laddr:  net.Addr{IP: "127.0.0.1", Port: port},
		Status: status,
		Pid:    pid,
	}
}

type fakeLister struct {
	conns []net.ConnectionStat
	err   error
}

func (f *fakeLister) Connections() ([]net.ConnectionStat, error) {
	return f.conns, f.err
}

type fakeKiller struct {
	killed []int32
	uids   map[int32][]int32
	killErr error
}

func (f *fakeKiller) Kill(pid int32) error {
	if f.killErr != nil {
		return f.killErr
	}
	f.killed = append(f.killed, pid)
	return nil
}

func (f *fakeKiller) Uids(pid int32) ([]int32, error) {
	if uids, ok := f.uids[pid]; ok {
		return uids, nil
	}
	return []int32{0}, nil
}

func testLogger() *log.Logger {
	return log.NewWithOptions(io.Discard, log.Options{ReportTimestamp: false})
}

func TestGetPIDFromConnections(t *testing.T) {
	tests := []struct {
		name        string
		port        string
		conns       []net.ConnectionStat
		expectedPID int32
		expectError bool
	}{
		{
			name: "find listening process",
			port: "8080",
			conns: []net.ConnectionStat{
				mockConnection(8080, 1234, "LISTEN"),
				mockConnection(3000, 5678, "LISTEN"),
			},
			expectedPID: 1234,
		},
		{
			name: "port exists but not listening",
			port: "8080",
			conns: []net.ConnectionStat{
				mockConnection(8080, 1234, "ESTABLISHED"),
			},
			expectedPID: -1,
		},
		{
			name: "port not found",
			port: "9999",
			conns: []net.ConnectionStat{
				mockConnection(8080, 1234, "LISTEN"),
			},
			expectedPID: -1,
		},
		{
			name:        "invalid port",
			port:        "abc",
			expectError: true,
		},
		{
			name:        "empty port",
			port:        "",
			expectError: true,
		},
		{
			name:        "port with special characters",
			port:        "80!",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid, err := getPIDFromConnections(tt.conns, tt.port)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for port %q, got none", tt.port)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if pid != tt.expectedPID {
				t.Errorf("got PID %d, want %d", pid, tt.expectedPID)
			}
		})
	}
}

func TestCanKill(t *testing.T) {
	tests := []struct {
		name       string
		currentUID int
		processUID int
		force      bool
		expected   bool
	}{
		{"user owns process", 1000, 1000, false, true},
		{"user does not own process", 1000, 1001, false, false},
		{"root user", 0, 1000, false, true},
		{"force flag", 1000, 1001, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canKill(tt.currentUID, tt.processUID, tt.force)
			if got != tt.expected {
				t.Errorf("canKill(%d, %d, %v) = %v, want %v", tt.currentUID, tt.processUID, tt.force, got, tt.expected)
			}
		})
	}
}

func TestRun_KillsMatchingProcess(t *testing.T) {
	killer := &fakeKiller{uids: map[int32][]int32{1234: {1000}}}
	k := &Krampus{
		Ports: []string{"8080"},
		UID:   1000,
		Lister: &fakeLister{conns: []net.ConnectionStat{
			mockConnection(8080, 1234, "LISTEN"),
		}},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 1 || killer.killed[0] != 1234 {
		t.Errorf("expected PID 1234 killed, got %v", killer.killed)
	}
}

func TestRun_SkipsUnownedProcess(t *testing.T) {
	killer := &fakeKiller{uids: map[int32][]int32{1234: {0}}}
	k := &Krampus{
		Ports: []string{"8080"},
		UID:   1000,
		Lister: &fakeLister{conns: []net.ConnectionStat{
			mockConnection(8080, 1234, "LISTEN"),
		}},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 0 {
		t.Errorf("expected no kills, got %v", killer.killed)
	}
}

func TestRun_ForceKillsUnownedProcess(t *testing.T) {
	killer := &fakeKiller{uids: map[int32][]int32{1234: {0}}}
	k := &Krampus{
		Ports: []string{"8080"},
		Force: true,
		UID:   1000,
		Lister: &fakeLister{conns: []net.ConnectionStat{
			mockConnection(8080, 1234, "LISTEN"),
		}},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 1 || killer.killed[0] != 1234 {
		t.Errorf("expected PID 1234 killed, got %v", killer.killed)
	}
}

func TestRun_MultiplePorts(t *testing.T) {
	killer := &fakeKiller{uids: map[int32][]int32{1234: {1000}, 5678: {1000}}}
	k := &Krampus{
		Ports: []string{"8080", "3000"},
		UID:   1000,
		Lister: &fakeLister{conns: []net.ConnectionStat{
			mockConnection(8080, 1234, "LISTEN"),
			mockConnection(3000, 5678, "LISTEN"),
		}},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 2 {
		t.Errorf("expected 2 kills, got %v", killer.killed)
	}
}

func TestRun_NoProcessOnPort(t *testing.T) {
	killer := &fakeKiller{}
	k := &Krampus{
		Ports:  []string{"9999"},
		UID:    1000,
		Lister: &fakeLister{conns: []net.ConnectionStat{}},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 0 {
		t.Errorf("expected no kills, got %v", killer.killed)
	}
}

func TestRun_ConnectionError(t *testing.T) {
	killer := &fakeKiller{}
	k := &Krampus{
		Ports:  []string{"8080"},
		UID:    1000,
		Lister: &fakeLister{err: fmt.Errorf("connection error")},
		Killer: killer,
		Logger: testLogger(),
	}

	k.Run()

	if len(killer.killed) != 0 {
		t.Errorf("expected no kills, got %v", killer.killed)
	}
}

func TestPrintVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"valid semver", "1.2.3", "v1.2.3\n"},
		{"invalid semver", "invalid", "0.0.0\n"},
		{"empty version", "", "0.0.0\n"},
		{"semver with prerelease", "1.2.3-alpha.1", "v1.2.3-alpha.1\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Version = tt.version

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printVersion()

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)

			if buf.String() != tt.expected {
				t.Errorf("printVersion() output = %q; want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestHasNonEmptyArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected bool
	}{
		{"empty array", []string{}, false},
		{"all empty strings", []string{"", "", ""}, false},
		{"has non-empty", []string{"foo", "bar"}, true},
		{"mixed", []string{"foo", "", "bar"}, true},
		{"single non-empty", []string{"test"}, true},
		{"single empty", []string{""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasNonEmptyArgs(tt.input)
			if got != tt.expected {
				t.Errorf("hasNonEmptyArgs(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkGetPIDFromConnections(b *testing.B) {
	conns := make([]net.ConnectionStat, 100)
	for i := range conns {
		conns[i] = mockConnection(uint32(8000+i), int32(1000+i), "LISTEN")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getPIDFromConnections(conns, "8050")
	}
}
