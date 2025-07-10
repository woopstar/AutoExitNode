package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/getlantern/systray"
)

var lastSSID string
var lastCellular bool
var lastCommand string
var tailscaleAvailable = true
var checkInterval = 15 * time.Second

// Use crypto/rand for secure random interval
var updateInterval = getSecureRandomInterval()

// getSecureRandomInterval returns a random duration between 1 and 60 minutes using crypto/rand.
func getSecureRandomInterval() time.Duration {
	n, err := rand.Int(rand.Reader, big.NewInt(60))
	if err != nil {
		// fallback to 15 minutes if crypto/rand fails
		return 15 * time.Minute
	}
	return time.Duration(n.Int64()+1) * time.Minute
}

func autoExitNote() {
	mStatus := systray.AddMenuItem("Status: Initializing...", "Current network status")
	mStatus.Disable()

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
			time.Sleep(checkInterval)
		}
	}()

	go func() {
		for {
			checkForUpdate(func(ver, url string) {
				updateVersionMenu(mVersion, ver, url)
			})
			time.Sleep(updateInterval)
		}
	}()
}

// checkAndApply handles the main logic for tray status and tailscale actions.
func checkAndApply(mStatus *systray.MenuItem) {
	if !tailscaleAvailable {
		return
	}
	ssid, err := getCurrentSSID()
	cell := isCellularConnected()
	online := hasInternetConnection()

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
	case !online:
		statusText = "No Internet"
		tooltip = fmt.Sprintf("Active: %s (unknown network)", getExitNodeName())
		icon = iconActive
		command = "activated"
		activateExitNode()
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

	mStatus.SetTitle("Status: " + statusText)
	systray.SetIcon(icon)
	systray.SetTooltip(tooltip)

	if ssid == lastSSID && cell == lastCellular && lastCommand == command && online {
		return
	}
	lastSSID = ssid
	lastCellular = cell

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

// updateVersionMenu updates the version menu item if a new version is available.
func updateVersionMenu(mVersion *systray.MenuItem, ver, url string) {
	if ver != "" && ver != currentVersion {
		mVersion.SetTitle(fmt.Sprintf("Version: %s (Update: %s)", currentVersion, ver))
		mVersion.SetTooltip(fmt.Sprintf("New version available: %s\n%s", ver, url))
	} else {
		mVersion.SetTitle(fmt.Sprintf("Version: %s", currentVersion))
		mVersion.SetTooltip("Current version")
	}
}
