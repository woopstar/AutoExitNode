package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

//go:embed icon_active.ico
var iconActive []byte

//go:embed icon_inactive.ico
var iconInactive []byte

type Config struct {
	TrustedSSIDs []string `json:"trustedSSIDs"`
	ExitNodes    []string `json:"exitNodes"`
}

var config Config
var tailscalePath = "C:\\Program Files\\Tailscale\\tailscale.exe"

var lastSSID string
var lastCellular bool
var lastCommand string // "activated" or "deactivated"
var tailscaleAvailable = true

var currentVersion = "v1.2.1" // Default version, will be overwritten by config if present

var latestVersion string
var latestVersionURL string

func main() {
	loadConfig()
	tailscaleAvailable = checkTailscaleExists()
	systray.Run(onReady, nil)
}

func loadConfig() {
	// Default config
	config = Config{
		TrustedSSIDs: []string{"Yoda-Fi", "R2D2-Fi"},
		ExitNodes:    []string{"homeassistant", "router", "vpn-node"},
	}
	f, err := os.Open("config.json")
	if err == nil {
		defer f.Close()
		_ = json.NewDecoder(f).Decode(&config)
	}
}

func checkTailscaleExists() bool {
	_, err := os.Stat(tailscalePath)
	return err == nil
}

func onReady() {
	// Status label
	mStatus := systray.AddMenuItem("Status: Initializing...", "Current network status")
	mStatus.Disable()

	// Version label
	mVersion := systray.AddMenuItem(fmt.Sprintf("Version: %s", currentVersion), "Current version")
	mVersion.Disable()

	mForce := systray.AddMenuItem("Force Sync", "Force immediate sync")
	mRunAtStartup := systray.AddMenuItemCheckbox("Run at startup", "Toggle auto-start", isStartupEnabled())
	mCheckUpdate := systray.AddMenuItem("Check for update", "Check for new version")
	mQuit := systray.AddMenuItem("Quit", "Exit the app")

	if !tailscaleAvailable {
		mForce.Disable()
		systray.SetIcon(iconInactive)
		systray.SetTooltip("Tailscale not found! Please install Tailscale.")
	} else {
		systray.SetIcon(iconInactive)
		systray.SetTooltip("AutoExitNode - Tailscale controller")
	}

	go func() {
		for {
			select {
			case <-mForce.ClickedCh:
				checkAndApply(mStatus)
			case <-mRunAtStartup.ClickedCh:
				if isStartupEnabled() {
					removeStartupShortcut()
					mRunAtStartup.Uncheck()
				} else {
					addStartupShortcut()
					mRunAtStartup.Check()
				}
			case <-mCheckUpdate.ClickedCh:
				go func() {
					checkForUpdate(func(ver, url string) {
						updateVersionMenu(mVersion, ver, url)
					})
				}()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	go func() {
		for {
			checkAndApply(mStatus)
			time.Sleep(15 * time.Second)
		}
	}()

	// Periodically check for updates in the background
	go func() {
		for {
			checkForUpdate(func(ver, url string) {
				updateVersionMenu(mVersion, ver, url)
			})
			time.Sleep(15 * time.Minute)
		}
	}()

	// Initial update check at startup
	go checkForUpdate(func(ver, url string) {
		updateVersionMenu(mVersion, ver, url)
	})
}

// Update the version menu item if a new version is available
func updateVersionMenu(mVersion *systray.MenuItem, ver, url string) {
	if ver != "" && ver != currentVersion {
		mVersion.SetTitle(fmt.Sprintf("Version: %s (Update: %s)", currentVersion, ver))
		mVersion.SetTooltip(fmt.Sprintf("New version available: %s\n%s", ver, url))
	} else {
		mVersion.SetTitle(fmt.Sprintf("Version: %s", currentVersion))
		mVersion.SetTooltip("Current version")
	}
}

func checkAndApply(mStatus *systray.MenuItem) {
	if !tailscaleAvailable {
		return
	}
	ssid, err := getCurrentSSID()
	cell := isCellularConnected()

	// Status label logic
	statusText := ""
	tooltip := ""
	var icon []byte = iconInactive
	command := ""

	switch {
	case cell:
		statusText = "Cellular"
		tooltip = fmt.Sprintf("Active: %s via cellular", getExitNodeName())
		icon = iconActive
		command = "activated"
	case err != nil || ssid == "":
		statusText = "Untrusted SSID"
		tooltip = fmt.Sprintf("Active: %s (unknown network)", getExitNodeName())
		icon = iconActive
		command = "activated"
	case isSSIDTrusted(ssid):
		statusText = fmt.Sprintf("Trusted SSID: %s", ssid)
		tooltip = fmt.Sprintf("Inactive: trusted network (%s)", ssid)
		icon = iconInactive
		command = "deactivated"
	default:
		statusText = "Untrusted SSID"
		tooltip = fmt.Sprintf("Active: %s (untrusted SSID)", getExitNodeName())
		icon = iconActive
		command = "activated"
	}

	// Update tray label, icon, tooltip
	mStatus.SetTitle("Status: " + statusText)
	systray.SetIcon(icon)
	systray.SetTooltip(tooltip)

	// Rate limiting: skip if nothing changed
	if ssid == lastSSID && cell == lastCellular && lastCommand == command {
		return
	}
	lastSSID = ssid
	lastCellular = cell

	// Only run tailscale if command changed
	if lastCommand != command {
		if command == "activated" {
			fmt.Println("[Activate] via", statusText)
			activateExitNode()
		} else {
			fmt.Println("[Deactivate] via", statusText)
			deactivateExitNode()
		}
		lastCommand = command
	}
}

func getExitNodeName() string {
	// Try each exit node in config, return first that works (for now, just return first)
	if len(config.ExitNodes) > 0 {
		return config.ExitNodes[0]
	}
	return "homeassistant"
}

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

func isSSIDTrusted(ssid string) bool {
	for _, trusted := range config.TrustedSSIDs {
		if strings.EqualFold(ssid, trusted) {
			return true
		}
	}
	return false
}

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

func getStartupShortcutPath() string {
	startupDir := os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup`
	return filepath.Join(startupDir, "AutoExitNode.lnk")
}

func isStartupEnabled() bool {
	_, err := os.Stat(getStartupShortcutPath())
	return err == nil
}

func addStartupShortcut() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	shellObj, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("CreateObject error:", err)
		return
	}
	defer shellObj.Release()

	wshell, err := shellObj.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("QueryInterface error:", err)
		return
	}
	defer wshell.Release()

	shortcutPath := getStartupShortcutPath()
	sc, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcutPath)
	if err != nil {
		fmt.Println("CreateShortcut error:", err)
		return
	}
	defer sc.Clear()

	shortcut := sc.ToIDispatch()
	_, _ = oleutil.PutProperty(shortcut, "TargetPath", exePath)
	_, _ = oleutil.PutProperty(shortcut, "WorkingDirectory", filepath.Dir(exePath))
	_, _ = oleutil.PutProperty(shortcut, "WindowStyle", 7)
	_, _ = oleutil.CallMethod(shortcut, "Save")
}

func removeStartupShortcut() {
	path := getStartupShortcutPath()
	if err := os.Remove(path); err != nil {
		fmt.Println("Failed to remove startup shortcut:", err)
	}
}

func checkForUpdate(cb func(version, url string)) {
	const repo = "woopstar/AutoExitNode" // Ensure this matches your GitHub repo (owner/repo)
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
		latestVersion = data.TagName
		latestVersionURL = data.HTMLURL
		showWindowsNotification("Update available!", fmt.Sprintf("New version: %s\nSee: %s", data.TagName, data.HTMLURL))
		if cb != nil {
			cb(data.TagName, data.HTMLURL)
		}
	} else if cb != nil {
		cb("", "")
	}
}

// showWindowsNotification displays a notification on Windows using go-toast.
func showWindowsNotification(title, message string) {
	// Show a simple popup using Windows MessageBox via PowerShell for maximum compatibility.
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
