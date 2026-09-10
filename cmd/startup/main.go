package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const defaultPort = "7006"
const serverCommand = "./cmd/server"

func main() {
	port := os.Getenv("OPENFEED_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := ensurePortAvailable(port); err != nil {
		fmt.Fprintf(os.Stderr, "OpenFeed startup blocked: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "run", serverCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "OpenFeed server exited: %v\n", err)
		os.Exit(1)
	}
}

func ensurePortAvailable(port string) error {
	if !portInUse(port) {
		return nil
	}

	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine repository root: %w", err)
	}

	processes, err := findOwnedOpenFeedProcesses(root)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf("port %s is already in use by another or unrecognized process; refusing to kill it", port)
	}

	for _, pid := range processes {
		if err := terminateProcess(pid); err != nil {
			return fmt.Errorf("failed to terminate previous OpenFeed process PID %d: %w", pid, err)
		}
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !portInUse(port) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("previous OpenFeed process exited but port %s was not released", port)
}

func portInUse(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 150*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

type processInfo struct {
	pid int
	cmd string
	cwd string
}

func findOwnedOpenFeedProcesses(root string) ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("cannot inspect process table: %w", err)
	}

	root, err = filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("cannot normalize repository root: %w", err)
	}

	var pids []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid == os.Getpid() {
			continue
		}

		info, err := readProcessInfo(pid)
		if err != nil {
			continue
		}
		if isManagedOpenFeedProcess(root, info) {
			pids = append(pids, info.pid)
		}
	}
	return pids, nil
}

func isManagedOpenFeedProcess(root string, info processInfo) bool {
	if isNativeOpenFeedCommand(info.cmd) {
		return true
	}
	return info.cwd == root && isOpenFeedCommand(info.cmd)
}

func readProcessInfo(pid int) (processInfo, error) {
	cmdline, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil {
		return processInfo{}, err
	}
	cmd := strings.Join(strings.FieldsFunc(string(cmdline), func(r rune) bool { return r == 0 }), " ")
	cwd, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(pid), "cwd"))
	if err != nil {
		return processInfo{}, err
	}
	return processInfo{pid: pid, cmd: cmd, cwd: cwd}, nil
}

func isNativeOpenFeedCommand(cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return false
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}
	return strings.EqualFold(filepath.Base(parts[0]), "openfeed")
}

func isOpenFeedCommand(cmd string) bool {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	return isNativeOpenFeedCommand(cmd) ||
		strings.Contains(cmd, "go run ./cmd/server") ||
		(strings.Contains(cmd, "go-build") && strings.HasSuffix(cmd, "/server")) ||
		(strings.Contains(cmd, "openfeed") && strings.Contains(cmd, "server"))
}

func terminateProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if !processExists(pid) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err := process.Signal(syscall.SIGKILL); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	return nil
}

func processExists(pid int) bool {
	file, err := os.Open(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return false
	}
	defer file.Close()
	return bufio.NewScanner(file).Scan()
}
