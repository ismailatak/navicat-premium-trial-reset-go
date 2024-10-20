package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const currentVersion = "v0.0.9"
const repoURL = "https://api.github.com/repos/ismailatak/navicat-premium-trial-reset-go/releases/latest"

func printVersion() {
	fmt.Println("Navicat Premium Trial Reset version", currentVersion)
}

// Fetch the latest version from GitHub
func checkForUpdate() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(repoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release info, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	latestVersion := result["tag_name"].(string)
	return latestVersion, nil
}

// Compare local version with the latest version
func ensureLatestVersion() {
	latestVersion, err := checkForUpdate()
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		os.Exit(1)
	}

	if latestVersion != currentVersion {
		fmt.Printf("A new version (%s) is available. Please run:\n", latestVersion)
		newGoInstall := fmt.Sprintf("  go install github.com/ismailatak/navicat-premium-trial-reset-go@%s", latestVersion)
		fmt.Println(newGoInstall)
		os.Exit(0)
	}
}
