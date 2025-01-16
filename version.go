package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const currentVersion = "v0.1.0"
const releasesURL = "https://api.github.com/repos/ismailatak/navicat-premium-trial-reset-go/releases"
const releaseLatestURL = releasesURL + "/latest"

func printVersion() {
	fmt.Println("Navicat Premium Trial Reset version", currentVersion)
}

// Fetch the releases from GitHub
func fetchReleases() ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", releasesURL, nil)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" { // 403 issue: Token required for higher rate limit on requests from GitHub Actions.
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases, status code: %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Fetch the latest release from GitHub
func checkForUpdate() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", releaseLatestURL, nil)
	if err != nil {
		return "", err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" { // 403 issue: Token required for higher rate limit on requests from GitHub Actions.
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release, status code: %d", resp.StatusCode)
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
	releases, err := fetchReleases()
	if err != nil {
		fmt.Println("Error fetching releases:", err)
		os.Exit(1)
	}

	if len(releases) == 0 {
		fmt.Println("No releases found")
		os.Exit(0)
	}

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
