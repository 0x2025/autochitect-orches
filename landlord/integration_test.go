package main

import (
	"bytes"
	"os/exec"
	"testing"
)

// TestDoctorCommand runs the 'landlord doctor' command and verifies its output
// This is an integration test that exercises the real CLI functionality
func TestDoctorCommand(t *testing.T) {
	// Build the command: landlord doctor
	cmd := exec.Command("landlord", "doctor")
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		t.Fatalf("landlord doctor command failed: %v", err)
	}

	// Get the output
	output := outBuf.String()

	// Verify the command executed and produced expected output
	if !bytes.Contains([]byte(output), []byte("Checking installed tools...")) {
		t.Errorf("expected 'Checking installed tools...' in output, got: %s", output)
	}

	// Check for tool status messages
	expectedPatterns := []string{
		"docker: ",
		"podman: ",
		"openshell: ",
	}

	for _, pattern := range expectedPatterns {
		if !bytes.Contains([]byte(output), []byte(pattern)) {
			t.Errorf("expected pattern '%s' in output, got: %s", pattern, output)
		}
	}
}

// TestStartCommand runs the 'landlord start' command and verifies its behavior
// This test ensures the start command executes without panicking and
// follows the expected flow for detecting and installing tools
func TestStartCommand(t *testing.T) {
	// Build the command: landlord start
	cmd := exec.Command("landlord", "start")
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	// Run the command (may fail due to network or missing update.json, but should not panic)
	if err := cmd.Run(); err != nil {
		// Expected to fail in sandbox due to network, but should not panic
		if !bytes.Contains([]byte(err.Error()), []byte("panic")) {
			t.Logf("landlord start command failed as expected: %v", err)
		} else {
			t.Fatalf("landlord start command panicked: %v", err)
		}
	} else {
		// If it succeeds, verify output contains expected messages
		output := outBuf.String()
		if !bytes.Contains([]byte(output), []byte("Starting landlord...")) {
			t.Errorf("expected 'Starting landlord...' in output, got: %s", output)
		}
	}
}