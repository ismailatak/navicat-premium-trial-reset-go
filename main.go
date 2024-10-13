package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	// Detect Navicat Premium version
	cmd := exec.Command("defaults", "read", "/Applications/Navicat Premium.app/Contents/Info.plist")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error reading Info.plist:", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`CFBundleShortVersionString = "([^\.]+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("Version not found")
		os.Exit(1)
	}

	version := matches[1]
	fmt.Printf("Detected Navicat Premium version %s\n", version)

	var file string

	// Select plist file by version
	switch version {
	case "17", "16":
		file = os.Getenv("HOME") + "/Library/Preferences/com.navicat.NavicatPremium.plist"
	case "15":
		file = os.Getenv("HOME") + "/Library/Preferences/com.prect.NavicatPremium15.plist"
	default:
		fmt.Printf("Version '%s' not handled\n", version)
		os.Exit(1)
	}

	fmt.Println("Resetting trial time...")

	// Delete hash from plist file
	cmd = exec.Command("defaults", "read", file)
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error reading plist file:", err)
		os.Exit(1)
	}

	re = regexp.MustCompile(`([0-9A-Z]{32}) = `)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		hash := matches[1]
		fmt.Printf("deleting %s array...\n", hash)
		exec.Command("defaults", "delete", file, hash).Run()
	}

	// Delete hidden file in Application Support
	cmd = exec.Command("ls", "-a", os.Getenv("HOME")+"/Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium/")
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	re = regexp.MustCompile(`\.([0-9A-Z]{32})`)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		hash2 := matches[1]
		fmt.Printf("deleting %s folder...\n", hash2)
		exec.Command("rm", os.Getenv("HOME")+"/Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium/."+hash2).Run()
	}

	fmt.Println("Done")
}
