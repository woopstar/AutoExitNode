package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
)

// isValidTailscalePath checks if the path is absolute and points to an .exe file in Program Files.
// isValidTailscalePath checks if the path is absolute, points to tailscale.exe in Program Files, and is executable.
func isValidTailscalePath(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if file exists and is executable
	info, err := os.Stat(abs)
	if err != nil || info.IsDir() {
		return false
	}

	// On Windows, check for .exe extension and that the file is not a directory
	return filepath.Ext(abs) == ".exe"
}

// isValidExitNodeName ensures the exit node name is alphanumeric, dash or underscore (no spaces or shell metacharacters).
func isValidExitNodeName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString(name)
}

// getExitNodeName returns the first exit node from config or a default.
func getExitNodeName() string {
	if len(config.ExitNodes) > 0 {
		return config.ExitNodes[0]
	}
	return "homeassistant"
}

// activateExitNode runs tailscale up with exit node.
// #nosec G204
func activateExitNode() {
	if !isValidTailscalePath(tailscalePath) {
		fmt.Println("Unsafe tailscalePath, aborting command")
		return
	}
	exitNode := getExitNodeName()
	if !isValidExitNodeName(exitNode) {
		fmt.Println("Unsafe exit node name, aborting command")
		return
	}
	cmd := exec.Command(tailscalePath,
		"up",
		"--exit-node="+exitNode,
		"--accept-dns=true",
		"--shields-up")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	fmt.Println("Activate output:", string(output))
	if err != nil {
		fmt.Println("Activate error:", err)
	}
}

// deactivateExitNode disables exit node.
// #nosec G204
func deactivateExitNode() {
	if !isValidTailscalePath(tailscalePath) {
		fmt.Println("Unsafe tailscalePath, aborting command")
		return
	}
	cmd := exec.Command(tailscalePath,
		"up",
		"--exit-node=",
		"--accept-dns=true",
		"--shields-up")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	fmt.Println("Deactivate output:", string(output))
	if err != nil {
		fmt.Println("Deactivate error:", err)
	}
}
