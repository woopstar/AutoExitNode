package main

import (
	_ "embed"
	"encoding/json"
	"os"

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
var currentVersion = "v1.3.0"
var startupDir = os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup`

func main() {
	loadConfig()
	tailscaleAvailable = checkTailscaleExists()
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
