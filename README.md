# AutoExitNode

[![GitHub Release][releases-shield]][releases]
[![GitHub Downloads][downloads-shield]][downloads]
[![License][license-shield]][license]
[![BuyMeCoffee][buymecoffeebadge]][buymecoffee]

![Icon](icon_active.png)

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

[releases-shield]: https://img.shields.io/github/v/release/woopstar/AutoExitNode?style=for-the-badge
[releases]: https://github.com/woopstar/AutoExitNode/releases
[downloads-shield]: https://img.shields.io/github/downloads/woopstar/AutoExitNode/total.svg?style=for-the-badge
[downloads]: https://github.com/woopstar/AutoExitNode/releases
[license-shield]: https://img.shields.io/github/license/woopstar/AutoExitNode?style=for-the-badge
[license]: https://github.com/woopstar/AutoExitNode/blob/main/LICENSE
[buymecoffeebadge]: https://img.shields.io/badge/buy%20me%20a%20coffee-donate-FFDD00.svg?style=for-the-badge&logo=buymeacoffee
[buymecoffee]: https://www.buymeacoffee.com/woopstar
