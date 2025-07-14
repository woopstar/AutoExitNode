package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/sys/windows"

	"github.com/getlantern/systray"
)

//go:embed assets/icon_active.ico
var iconActive []byte

//go:embed assets/icon_inactive.ico
var iconInactive []byte

type Config struct {
	TailscalePath string   `json:"tailscalePath"`
	TrustedSSIDs  []string `json:"trustedSSIDs"`
	ExitNodes     []string `json:"exitNodes"`
}

var config Config
var tailscalePath string
var currentVersion = "v2.1.2"
var startupDir = os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup`
var mutexHandle windows.Handle

func main() {
	mutexName, _ := windows.UTF16PtrFromString("Global\\AutoExitNodeMutex")
	handle, err := windows.CreateMutex(nil, false, mutexName)
	if err != nil {
		fmt.Println("Failed to create mutex:", err)
		os.Exit(1)
	}
	mutexHandle = handle
	lastErr := windows.GetLastError()
	if lastErr == windows.ERROR_ALREADY_EXISTS {
		fmt.Println("Only one instance is allowed.")
		os.Exit(0)
	}
	defer windows.ReleaseMutex(mutexHandle)
	defer windows.CloseHandle(mutexHandle)

	loadConfig()
	tailscaleAvailable = isValidTailscalePath(tailscalePath)
	systray.Run(autoExitNote, nil)
}

func loadConfig() {
	// Default config
	config = Config{
		TailscalePath: "C:\\Program Files\\Tailscale\\tailscale.exe",
		TrustedSSIDs:  []string{"Yoda-Fi", "R2D2-Fi"},
		ExitNodes:     []string{"homeassistant", "router", "vpn-node"},
	}
	f, err := os.Open("config.json")
	if err == nil {
		defer f.Close()
		_ = json.NewDecoder(f).Decode(&config)
	}

	tailscalePath = config.TailscalePath
}
