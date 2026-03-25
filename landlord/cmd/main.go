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
	fetchUpdate()
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

func fetchUpdate() {
	resp, err := http.Get("https://autochitect.com/landlord/update.json")
	if err != nil {
		fmt.Printf("Failed to fetch update.json: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var update map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
		fmt.Printf("Failed to parse update.json: %v\n", err)
		return
	}

	// Handle OS-specific installation
	osType := "linux"
	if runtime.GOOS == "darwin" {
		osType = "macos"
	}

	// Install missing tools based on OS
	if !isInstalled("docker") {
		fmt.Println("Installing docker...")
		installTool("docker", getToolURL(update, "docker", osType))
	}
	if !isInstalled("podman") {
		fmt.Println("Installing podman...")
		installTool("podman", getToolURL(update, "podman", osType))
	}
	if !isInstalled("openshell") {
		fmt.Println("Installing openshell...")
		installTool("openshell", getToolURL(update, "openshell", osType))
	}

	// Report OpenShell version from server
	if version, ok := update["openshell_version"]; ok {
		fmt.Printf("OpenShell version from server: %v\n", version)
	} else {
		fmt.Println("No OpenShell version info in update.json")
	}
}

func getToolURL(update map[string]interface{}, tool, osType string) string {
	// Get URL from update.json based on tool and OS
	switch tool {
	case "docker":
		return getDockerURL(update, osType)
	case "podman":
		return getPodmanURL(update, osType)
	case "openshell":
		return getOpenShellURL(update, osType)
	default:
		return ""
	}
}

func getDockerURL(update map[string]interface{}, osType string) string {
	// Docker installation URLs based on OS
	switch osType {
	case "macos":
		return "https://desktop.docker.com/mac/stable/Docker.dmg"
	case "linux":
		return "https://download.docker.com/linux/ubuntu/dists/" +
			"$(lsb_release -cs)/pool/stable/x86_64/docker-ce-cli_" +
			"$(apt-cache madison docker-ce-cli | head -1 | awk '{print $3}')_amd64.deb"
	default:
		return ""
	}
}

func getPodmanURL(update map[string]interface{}, osType string) string {
	// Podman installation URLs based on OS
	switch osType {
	case "macos":
		return "https://github.com/containers/podman/releases/download/latest/podman-desktop-mac.tar.gz"
	case "linux":
		return "https://github.com/containers/podman/releases/download/latest/podman-latest.tar.gz"
	default:
		return ""
	}
}

func getOpenShellURL(update map[string]interface{}, osType string) string {
	// OpenShell installation URLs based on OS
	switch osType {
	case "macos":
		return update["openshell_url"].(string)
	case "linux":
		return update["openshell_url"].(string)
	default:
		return ""
	}
}

func installTool(name, url string) error {
	fmt.Printf("Downloading %s...\n", name)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", name, err)
	}
	defer resp.Body.Close()

	// Save to temporary file
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Copy response to file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}
	tmpFile.Close()

	// Install based on OS and tool type
	switch name {
	case "docker":
		if runtime.GOOS == "darwin" {
			// For macOS: mount and copy app
			fmt.Println("Mounting Docker for macOS...")
			// Actual mounting logic would go here
		} else {
			// For Linux: install Debian package
			fmt.Println("Installing Docker for Linux...")
			// Actual installation logic would go here
		}
	case "podman":
		if runtime.GOOS == "darwin" {
			// For macOS: extract and install
			fmt.Println("Extracting Podman for macOS...")
			// Actual extraction logic would go here
		} else {
			// For Linux: extract and move binaries
			fmt.Println("Extracting Podman for Linux...")
			// Actual extraction logic would go here
		}
	case "openshell":
		if runtime.GOOS == "darwin" {
			// For macOS: mount and copy app
			fmt.Println("Mounting OpenShell for macOS...")
			// Actual mounting logic would go here
		} else {
			// For Linux: extract and move binaries
			fmt.Println("Extracting OpenShell for Linux...")
			// Actual extraction logic would go here
		}
	default:
		return fmt.Errorf("unknown tool: %s", name)
	}

	// Verify installation
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s installation failed: %w", name, err)
	}
	fmt.Printf("%s installed successfully\n", name)
	return nil
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}