package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// checkForUpdate checks for a new release and calls cb if found.
func checkForUpdate(cb func(version, url string)) {
	const repo = "woopstar/AutoExitNode"
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("checkForUpdate: failed to create request:", err)
		return
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("checkForUpdate: HTTP request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("checkForUpdate: unexpected status code: %d\n", resp.StatusCode)
		return
	}

	var data struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("checkForUpdate: failed to decode JSON:", err)
		return
	}

	if data.TagName != "" && data.TagName != currentVersion {
		showWindowsNotification("Update available!", fmt.Sprintf("New version: %s\nSee: %s", data.TagName, data.HTMLURL))
		if cb != nil {
			cb(data.TagName, data.HTMLURL)
		}
	} else if cb != nil {
		cb("", "")
	}
}

// showWindowsNotification displays a notification on Windows using PowerShell.
func showWindowsNotification(title, message string) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework;[System.Windows.MessageBox]::Show('%s', '%s')`, escapeForPowerShell(message), escapeForPowerShell(title)))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		fmt.Println("showWindowsNotification: failed to show popup:", err)
	}
}

// escapeForPowerShell escapes single quotes for PowerShell string literals.
func escapeForPowerShell(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
