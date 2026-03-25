package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read update.json: %v\n", err)
		return
	}
	var update map[string]interface{}
	if err := json.Unmarshal(body, &update); err != nil {
		fmt.Printf("Failed to parse update.json: %v\n", err)
		return
	}
	if version, ok := update["openshell_version"]; ok {
		fmt.Printf("OpenShell version from server: %v\n", version)
	} else {
		fmt.Println("No OpenShell version info in update.json")
	}
}
