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

// main is the entry point of the landlord CLI application.
// It parses command-line arguments and dispatches to the appropriate command handler.
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

// doctor checks the installation status of docker, podman, and openshell.
// It prints whether each tool is installed and its version if available.
func doctor() {
	fmt.Println("Checking installed tools...")
	checkTool("docker")
	checkTool("podman")
	checkTool("openshell")
}

// start initiates the landlord start process.
// It checks for required tools, fetches the update configuration,
// and installs any missing tools based on the OS and configuration.
func start() {
	fmt.Println("Starting landlord...")
	checkTool("docker")
	checkTool("openshell")
	fetchAndInstall()
}

// checkTool verifies if a given command is available in the system PATH.
// If not installed, it reports that the tool is missing.
// If installed, it attempts to retrieve its version and prints the result.
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

// fetchAndInstall handles the installation of missing tools.
// It fetches the update configuration from the remote server,
// determines the current OS, and installs any missing tools using URLs
// specified in the configuration. It provides user-friendly messages
// for success, failures, and network issues.
func fetchAndInstall() {
	// Fetch update.json configuration
	resp, err := http.Get("https://autochitect.com/landlord/update.json")
	if err != nil {
		fmt.Printf("Network error fetching update.json: %v. Please check your connection and try again.\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read update.json: %v\n", err)
		return
	}

	var update map[string]interface{}
	if err := json.Unmarshal(data, &update); err != nil {
		fmt.Printf("Failed to parse update.json: %v\n", err)
		return
	}

	// Determine OS type
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
				fmt.Printf("Could not determine download URL for %s: %v\n", tool, err)
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
		fmt.Println("No OpenShell version information available in update.json")
	}
}

// isInstalled checks if a command exists in the system PATH.
func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// getToolURL retrieves the download URL for a specific tool and OS type
// from the update configuration. It returns an error if the tool or URL
// is not defined in the configuration.
func getToolURL(update map[string]interface{}, tool, osType string) (string, error) {
	toolData, ok := update[tool].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("%s configuration not found", tool)
	}
	platformData, ok := toolData[osType].(string)
	if !ok {
		return "", fmt.Errorf("URL for %s on %s not found", tool, osType)
	}
	return platformData.(string), nil
}

// downloadAndInstall downloads the specified tool from the given URL
// and installs it. This is a placeholder implementation that should be
// replaced with actual installation logic for production use.
// It returns an error if the download or installation fails.
func downloadAndInstall(tool, url string) error {
	fmt.Printf("Downloading %s from %s...\n", tool, url)
	// Placeholder: simulate download and installation
	// In a real implementation, this would use curl or wget to fetch and install the tool
	return nil
}