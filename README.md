# Navicat Premium Trial Reset (Go)

This Go script helps you reset the trial period of **Navicat Premium** by cleaning up license-related files and configurations. It supports Navicat versions **15**, **16**, and **17**.

## Features
- Automatically detects the installed Navicat Premium version.
- Cleans up trial-related plist and hidden folder files.
- Simple and fast execution using Go's `os/exec` package.
- Supports Navicat Premium 15, 16, and 17.

## Prerequisites
- **Go** installed on your system. You can download and install it from [here](https://go.dev/dl/).
- A MacOS system (this script relies on MacOS-specific commands like `defaults`).

## Installation

### Option 1: Using `go get`

1. Clone the repository:
   ```bash
   git clone https://github.com/ismailatak/navicat-premium-trial-reset-go.git
   cd navicat-premium-trial-reset-go
   ```

2. Build the script:
   ```bash
   go build -o navicat-reset
   ```

3. Run the script:
   ```bash
   ./navicat-reset
   ```

### Option 2: Using `go install`

1. Alternatively, you can install and run the script with Go's install command:
   ```bash
   go install github.com/ismailatak/navicat-premium-trial-reset-go@latest
   ```

2. After this, you can run the script directly as:
   ```bash
   navicat-reset
   ```

## How It Works
- The script reads the Info.plist file of the installed Navicat Premium app to detect the current version.
- Based on the version, it targets the correct plist file in ~/Library/Preferences/.
- It also scans and deletes hidden license files from ~/Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium/.

## Example Output
   ```bash
   Detected Navicat Premium version 17
   Resetting trial time...
   deleting [HASH] array...
   deleting [HASH] folder...
   Done
   ```

## Disclaimer

This script is for educational purposes only. The author does not encourage or support software piracy or unauthorized usage of licensed software.