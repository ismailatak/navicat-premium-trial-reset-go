package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Flags
	versionFlag := flag.Bool("version", false, "Display the current version")
	flag.Parse()

	// Print version only
	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	// Check for updates
	ensureLatestVersion()

	// Navicat Premium
	cmd := exec.Command("defaults", "read", "/Applications/Navicat Premium.app/Contents/Info.plist")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Please make sure Navicat Premium is installed")
		fmt.Println("Error reading Navicat Premium:", err)
		os.Exit(1)
	}

	// Detect Navicat Premium version
	re := regexp.MustCompile(`CFBundleShortVersionString = "([^\"]+)"`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("Error detecting Navicat Premium version")
		os.Exit(1)
	}
	fullVersion := matches[1]
	fmt.Printf("Detected Navicat Premium version %s\n", fullVersion)

	// Check if Navicat Premium is running
	if isRunning("Navicat Premium") {
		fmt.Println("")
		fmt.Println("Navicat Premium is currently running!")
		fmt.Println("Please save your work before continuing")
		fmt.Println("")

		fmt.Print("Press [enter] to close Navicat Premium and continue...")
		fmt.Scanln()

		if isRunning("Navicat Premium") {
			fmt.Println("Closing Navicat Premium...")
			fmt.Println("")

			err = exec.Command("killall", "Navicat Premium").Run()
			if err != nil {
				fmt.Println("Error Navicat Premium could not be closed!", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("The user closed Navicat Premium. I'm continuing with the process...")
			fmt.Println("")
		}

		time.Sleep(1 * time.Second)
	}

	// Select service name by version
	var serviceName string
	majorVersion := strings.Split(fullVersion, ".")[0]
	switch majorVersion {
	case "17", "16":
		serviceName = "com.navicat.NavicatPremium"
	case "15":
		serviceName = "com.prect.NavicatPremium15"
	default:
		fmt.Printf("Unsupported Navicat Premium version: %s\n", fullVersion)
		os.Exit(1)
	}

	fmt.Println("Resetting trial time...")

	// Get home env
	home := os.Getenv("HOME")

	// Preferences .plist file exists
	preferencesPListFile := home + "/Library/Preferences/" + serviceName + ".plist"
	exists := exec.Command("ls", "-l", "-a", preferencesPListFile)
	output, err = exists.Output()
	if err != nil || len(output) == 0 {
		fmt.Println("Preferences .plist file does not exist or is empty:", err)
		os.Exit(1)
	}

	// Preferences .plist file reading
	cmd = exec.Command("defaults", "read", preferencesPListFile)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error reading preferences .plist file (%s): %+v\n", preferencesPListFile, err)
		os.Exit(1)
	}

	// Delete hash from preferences .plist file
	re = regexp.MustCompile(`([0-9A-Z]{32}) = `)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		preferencesPListHash := matches[1]
		fmt.Printf("deleting preferences .plist hash: %s\n", preferencesPListHash)
		err = exec.Command("defaults", "delete", preferencesPListFile, preferencesPListHash).Run()
		if err != nil {
			fmt.Println("Error deleting preferences .plist hash:", err)
			os.Exit(1)
		}
	}

	// Application Support file exists
	applicationSupportFile := home + "/Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium/"
	exists = exec.Command("ls", "-a", applicationSupportFile)
	output, err = exists.Output()
	if err != nil {
		fmt.Println("Application Support file does not exist or is empty:", err)
		os.Exit(1)
	}

	// Delete hash from application support file
	re = regexp.MustCompile(`\.([0-9A-Z]{32})`)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		applicationSupportHash := matches[1]
		fmt.Printf("deleting application support hash: %s\n", applicationSupportHash)
		err = exec.Command("rm", applicationSupportFile+"."+applicationSupportHash).Run()
		if err != nil {
			fmt.Println("Error deleting application support hash:", err)
			os.Exit(1)
		}
	}

	// Keychain inspection
	needsKeychain := false
	if majorVersion == "17" {
		parts := strings.Split(fullVersion, ".")
		if len(parts) >= 3 {
			var major, minor, patch int
			fmt.Sscanf(parts[0], "%d", &major)
			fmt.Sscanf(parts[1], "%d", &minor)
			fmt.Sscanf(parts[2], "%d", &patch)

			if minor > 3 || (minor == 3 && patch >= 7) {
				needsKeychain = true
			}
		}
	}

	// Keychain deletion
	if needsKeychain {
		keychains := home + "/Library/Keychains/login.keychain-db"
		cmd = exec.Command("security", "dump-keychain", keychains)
		output, err = cmd.Output()
		if err != nil {
			fmt.Printf("Error reading keychains file (%s): %+v\n", keychains, err)
			os.Exit(1)
		}
		keychainLines := strings.Split(string(output), "\n")

		var keychainBlocks []string
		for i, kcLine := range keychainLines {
			if strings.Contains(kcLine, serviceName) {
				end := min(i+6, len(keychainLines))
				keychainBlocks = append(keychainBlocks, strings.Join(keychainLines[i:end], "\n"))
			}
		}

		re = regexp.MustCompile(`[0-9A-F]{32}`)
		keychainsHash := make(map[string]struct{})
		for _, kcBlock := range keychainBlocks {
			for _, kcBlockLine := range strings.Split(kcBlock, "\n") {
				if strings.Contains(kcBlockLine, "acct") {
					acctMatch := re.FindString(kcBlockLine)
					if acctMatch != "" {
						keychainsHash[acctMatch] = struct{}{}
					}
				}
			}
		}

		for kcHash := range keychainsHash {
			fmt.Printf("deleting keychains hash: %s\n", kcHash)
			cmd = exec.Command("security", "delete-generic-password", "-s", serviceName, "-a", kcHash)

			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error deleting keychains hash:", err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("Done")
}

func isRunning(name string) bool {
	err := exec.Command("pgrep", "-x", name).Run()
	return err == nil
}
