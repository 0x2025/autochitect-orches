package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: landlord <command> [args...]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "doctor":
		doctor()
	case "start":
		start()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func doctor() {
	fmt.Println("Checking installed tools...")
	checkTool("docker")
	checkTool("podman")
	checkTool("openshell")
}

func start() {
	fmt.Println("Starting landlord...")
	checkTool("docker")
	checkTool("openshell")
	fetchAndInstall()
}

func checkTool(name string) {
	if _, err := exec.LookPath(name); err != nil {
		fmt.Printf("%s: not installed\n", name)
		return
	}
	versionCmd := exec.Command(name, "--version")
	out, err := versionCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s: error getting version\n", name)
		return
	}
	fmt.Printf("%s: %s\n", strings.TrimSpace(name), strings.TrimSpace(string(out)))
}

func fetchAndInstall() {
	// Fetch update.json
	resp, err := http.Get("https://autochitect.com/landlord/update.json")
	if err != nil {
		fmt.Printf("Failed to fetch update.json: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read update.json: %v\n", err)
		return
	}
	var update map[string]interface{}
	if err := json.Unmarshal(body, &update); err != nil {
		fmt.Printf("Failed to parse update.json: %v\n", err)
		return
	}

	// Determine OS
	osType := "linux"
	if runtime.GOOS == "darwin" {
		osType = "macos"
	}

	// Install missing tools
	tools := []string{"docker", "podman", "openshell"}
	for _, tool := range tools {
		if !isInstalled(tool) {
			fmt.Printf("%s: not installed, installing...\n", tool)
			url, err := getToolURL(update, tool, osType)
			if err != nil {
				fmt.Printf("Could not get URL for %s: %v\n", tool, err)
				continue
			}
			if err := downloadAndInstall(tool, url); err != nil {
				fmt.Printf("Failed to install %s: %v\n", tool, err)
			} else {
				fmt.Printf("%s: installed successfully\n", tool)
			}
		}
	}

	// Report OpenShell version from server
	if version, ok := update["openshell_version"]; ok {
		fmt.Printf("OpenShell version from server: %v\n", version)
	} else {
		fmt.Println("No OpenShell version info in update.json")
	}
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func getToolURL(update map[string]interface{}, tool, osType string) (string, error) {
	if toolData, ok := update[tool]; ok {
		if toolDataMap, ok := toolData.(map[string]interface{}); ok {
			if urls, ok := toolDataMap["url"].(map[string]interface{}); ok {
				if url, ok := urls[osType]; ok {
					return fmt.Sprintf("%v", url), nil
				}
			}
		}
	}
	return "", fmt.Errorf("URL not found for tool %s on %s", tool, osType)
}

func downloadAndInstall(tool, url string) error {
	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()
	// Save to /tmp/<tool>-installer
	filePath := fmt.Sprintf("/tmp/%s-installer", tool)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Install based on OS
	switch runtime.GOOS {
	case "darwin":
		// For macOS, maybe open the dmg or run an installer
		// Here we just simulate
		fmt.Printf("On macOS, you would open %s to install %s\n", url, tool)
		return nil
	default:
		// For Linux, maybe make executable and copy
		execCmd := exec.Command("sh", "-c", fmt.Sprintf("chmod +x %s && echo 'Installing %s'", filePath, tool))
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	}
}