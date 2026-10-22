package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	list := flag.Bool("l", false, "list process on port (do not kill)")
	force := flag.Bool("f", false, "force kill (SIGKILL / taskkill /F)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "portkill — free a TCP port by stopping the process that owns it\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  portkill [-l] [-f] <port>\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n  portkill 3000\n  portkill -l 8080\n  portkill -f 5173\n")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil || port < 1 || port > 65535 {
		fmt.Fprintln(os.Stderr, "error: port must be 1–65535")
		os.Exit(2)
	}

	pids, detail, err := findPIDs(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if len(pids) == 0 {
		fmt.Printf("nothing listening on :%d\n", port)
		return
	}

	fmt.Printf("port %d → %s\n", port, detail)
	if *list {
		return
	}

	for _, pid := range pids {
		if err := killPID(pid, *force); err != nil {
			fmt.Fprintf(os.Stderr, "failed to stop pid %d: %v\n", pid, err)
			os.Exit(1)
		}
		fmt.Printf("stopped pid %d\n", pid)
	}
}

func findPIDs(port int) ([]int, string, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		out, err := exec.Command("lsof", "-nP", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN", "-t").CombinedOutput()
		if err != nil {
			// lsof exits 1 when nothing found
			if len(strings.TrimSpace(string(out))) == 0 {
				return nil, "", nil
			}
			return nil, "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		var pids []int
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line == "" {
				continue
			}
			pid, err := strconv.Atoi(line)
			if err == nil {
				pids = append(pids, pid)
			}
		}
		detail := strings.TrimSpace(string(out))
		if name := procName(pids); name != "" {
			detail = name + " (pid " + detail + ")"
		}
		return unique(pids), detail, nil
	default:
		return nil, "", fmt.Errorf("unsupported OS: %s (use Linux/macOS)", runtime.GOOS)
	}
}

func procName(pids []int) string {
	if len(pids) == 0 {
		return ""
	}
	out, err := exec.Command("ps", "-p", strconv.Itoa(pids[0]), "-o", "comm=").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func killPID(pid int, force bool) error {
	sig := os.Interrupt
	if force {
		return exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func unique(in []int) []int {
	seen := map[int]bool{}
	var out []int
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
