package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// checkTailscaleExists returns true if the Tailscale binary exists.
func checkTailscaleExists() bool {
	_, err := os.Stat(tailscalePath)
	return err == nil
}

// getExitNodeName returns the first exit node from config or a default.
func getExitNodeName() string {
	if len(config.ExitNodes) > 0 {
		return config.ExitNodes[0]
	}
	return "homeassistant"
}

// activateExitNode runs tailscale up with exit node.
func activateExitNode() {
	cmd := exec.Command(tailscalePath,
		"up",
		"--exit-node="+getExitNodeName(),
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
func deactivateExitNode() {
	cmd := exec.Command(tailscalePath,
		"up",
		"--exit-node=",
		"--accept-dns=false",
		"--shields-up")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	fmt.Println("Deactivate output:", string(output))
	if err != nil {
		fmt.Println("Deactivate error:", err)
	}
}
