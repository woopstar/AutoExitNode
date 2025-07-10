package main

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// getCurrentSSID returns the current WiFi SSID or an error if not found.
func getCurrentSSID() (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^\s*SSID\s*:\s(.*)$`)
	for _, line := range strings.Split(string(output), "\n") {
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1]), nil
		}
	}
	return "", errors.New("SSID not found")
}

// isCellularConnected returns true if a cellular interface is connected.
func isCellularConnected() bool {
	cmd := exec.Command("netsh", "mbn", "show", "interfaces")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "state") && strings.Contains(strings.ToLower(line), "connected") {
			return true
		}
	}
	return false
}

// isSSIDTrusted checks if the given SSID is in the trusted list.
func isSSIDTrusted(ssid string) bool {
	for _, trusted := range config.TrustedSSIDs {
		if strings.EqualFold(ssid, trusted) {
			return true
		}
	}
	return false
}

// hasInternetConnection checks if the device has an active internet connection.
func hasInternetConnection() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	if err != nil {
		conn, err := net.DialTimeout("tcp", "1.1.1.1:53", timeout)
		if err != nil {
			if conn != nil {
				if cerr := conn.Close(); cerr != nil {
					// Log error if closing connection fails
					fmt.Println("Error closing connection:", cerr)
				}
			}
			return false
		}
	}
	if conn != nil {
		if cerr := conn.Close(); cerr != nil {
			fmt.Println("Error closing connection:", cerr)
		}
	}
	return true
}
