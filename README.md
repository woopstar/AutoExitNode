# AutoExitNode

**AutoExitNode** is a Windows system tray application that automatically manages your Tailscale exit node based on your network (WiFi SSID or cellular connection).

## Features

- Shows active status in the tray menu (Trusted SSID, Untrusted SSID, Cellular)
- Dynamic icon (blue = active, gray = inactive)
- Tooltip with status and active exit node
- Supports fallback exit node names via `config.json`
- Rate limiting: avoids repeated Tailscale calls if nothing has changed
- Error handling if Tailscale is not installed
- Self-update check against GitHub releases
- Menu to force sync, toggle autostart, check for updates, and quit the app
- Displays running version in the tray menu

## Installation

1. **Install Tailscale**
   Download and install Tailscale from [https://tailscale.com/download](https://tailscale.com/download).

2. **Download AutoExitNode**
   Get the latest release from [GitHub Releases](https://github.com/andreas-kruger/AutoExitNode/releases).

3. **Place the program and icons**
   Make sure the following files are in the same folder:
   - `AutoExitNode.exe`
   - `icon_active.ico`
   - `icon_inactive.ico`
   - `config.json` (optional, see below)

4. **(Optional) Edit config.json**
   Example `config.json`:
   ```json
   {
     "trustedSSIDs": ["Yoda-Fi", "R2D2-Fi"],
     "exitNodes": ["homeassistant", "router", "vpn-node"]
   }
   ```

## Usage

- The app starts in the system tray and automatically manages the Tailscale exit node:
  - **Trusted SSID:** Disables exit node
  - **Untrusted SSID/Cellular:** Enables exit node (first in the config list)
- The tray menu shows status, version, and provides access to:
  - Force Sync (trigger immediate update)
  - Run at startup (autostart)
  - Check for update (checks GitHub for new version)
  - Quit (exit the app)

## Development & Testing

- Run `go build` to build the program.
- Unit tests are in `main_test.go`:
  ```
  go test
  ```

## Requirements

- Windows 10 or newer
- [Tailscale](https://tailscale.com/)
- Go 1.18+ (to build from source)

## Icons

- `icon_active.ico` (blue): Shown when exit node is active
- `icon_inactive.ico` (gray): Shown when exit node is inactive

## Updates

The app automatically checks for new versions on GitHub and shows a Windows notification if an update is available.

---

**Note:**
This project is not officially affiliated with Tailscale.
