package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v3/client"
)

var gisturl = "https://gist.githubusercontent.com/w0ofix/a8bdea9a12680ccbe3694026b371f33f/raw/gistfile1.txt"
var firstUpdateCheck = true
var updateWarning = false

func version() string {
	return "aero-0:89a53d4b2860a317"
}

func isUpdated() (bool, struct {
	Version         string   `json:"version"`
	AllowedVersions []string `json:"allowed_versions"`
	On              bool     `json:"on"`
}) {
	var latestVersion struct {
		Version         string   `json:"version"`
		AllowedVersions []string `json:"allowed_versions"`
		On              bool     `json:"on"`
	}

	cl := client.New()

	resp, err := cl.Get(gisturl)
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return false, latestVersion
	}

	if resp.StatusCode() != 200 {
		fmt.Println("Error checking for updates: received status code", resp.StatusCode())
		return false, latestVersion
	}

	if err := resp.JSON(&latestVersion); err != nil {
		fmt.Println("Error parsing version response:", err)
		return false, latestVersion
	}

	return latestVersion.Version == version(), latestVersion
}

/* Public methods */

func Guard() {
	if firstUpdateCheck {
		fmt.Println("[*] Checking for updates...")
	}

	updated, data := isUpdated()
	if !updated {
		if !updateWarning {
			fmt.Println("[-] A new version of TLS Api is available. Please update to the latest version")
			fmt.Println("[*] Actual version : " + version() + " | Latest version : " + data.Version)

			updateWarning = true
		}

		isAllowed := false
		for _, v := range data.AllowedVersions {
			if v == version() {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			fmt.Println("[-] Your version is no longer supported. You need to update to the latest version")
			os.Exit(-1)
		}
	}

	if updated && firstUpdateCheck {
		fmt.Println("[+] You are using the latest version of TLS Api | Version : " + version())
	}

	if !data.On {
		fmt.Println("[-] API Error : cac8ce")
		os.Exit(-1)
	}

	firstUpdateCheck = false
}
func StartPeriodicUpdateCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			updated, data := isUpdated()
			if !updated {
				if !updateWarning {
					fmt.Println("[-] A new version of TLS Api is available. Please update to the latest version")
					fmt.Println("[*] Actual version : " + version() + " | Latest version : " + data.Version)

					updateWarning = true
				}

				isAllowed := false
				for _, v := range data.AllowedVersions {
					if v == version() {
						isAllowed = true
						break
					}
				}

				if !isAllowed {
					fmt.Println("[-] Your version is no longer supported. You need to update to the latest version")
					os.Exit(-1)
				}
			}

			if !data.On {
				fmt.Println("[-] API Error : cac8ce")
				os.Exit(-1)
			}
		}
	}()
}
