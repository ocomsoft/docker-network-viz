// Package integration provides CLI integration tests for docker-network-viz.
package integration

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_HelpCommand tests the --help flag output.
func TestCLI_HelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "..", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("help command failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()

	expectedElements := []string{
		"docker-network-viz",
		"visualizing Docker network topology",
		"--config",
		"--no-color",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("expected help output to contain %q, got:\n%s", expected, output)
		}
	}
}

// TestCLI_VisualizeHelpCommand tests the visualize --help flag output.
func TestCLI_VisualizeHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "..", "visualize", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("visualize help command failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()

	expectedElements := []string{
		"visualize",
		"Visualize Docker network topology",
		"--only-network",
		"--container",
		"--no-aliases",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("expected visualize help output to contain %q, got:\n%s", expected, output)
		}
	}
}

// TestCLI_UnknownFlag tests that unknown flags produce an error.
func TestCLI_UnknownFlag(t *testing.T) {
	cmd := exec.Command("go", "run", "..", "--unknown-flag")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("expected error for unknown flag, got none")
	}

	// The error message should indicate unknown flag
	combinedOutput := stdout.String() + stderr.String()
	if !strings.Contains(combinedOutput, "unknown") {
		t.Logf("output did not contain 'unknown': %s", combinedOutput)
	}
}

// TestCLI_InvalidConfigFile tests behavior with non-existent config file.
func TestCLI_InvalidConfigFile(t *testing.T) {
	cmd := exec.Command("go", "run", "..", "--config", "/nonexistent/config.yaml", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// This should still show help since config errors are handled gracefully
	err := cmd.Run()
	if err != nil {
		t.Fatalf("command with invalid config should still work with --help: %v\nstderr: %s", err, stderr.String())
	}
}

// TestCLI_NoColorFlag tests that --no-color flag is accepted.
func TestCLI_NoColorFlag(t *testing.T) {
	cmd := exec.Command("go", "run", "..", "--no-color", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("--no-color flag should be accepted: %v\nstderr: %s", err, stderr.String())
	}
}

// TestCLI_VisualizeFlags tests that visualize command flags are accepted.
func TestCLI_VisualizeFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "only-network flag with help",
			args: []string{"visualize", "--only-network", "bridge", "--help"},
		},
		{
			name: "container flag with help",
			args: []string{"visualize", "--container", "test", "--help"},
		},
		{
			name: "no-aliases flag with help",
			args: []string{"visualize", "--no-aliases", "--help"},
		},
		{
			name: "multiple flags with help",
			args: []string{"visualize", "--only-network", "bridge", "--no-aliases", "--help"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"run", ".."}, tc.args...)
			cmd := exec.Command("go", args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("command with flags %v should be accepted: %v\nstderr: %s", tc.args, err, stderr.String())
			}
		})
	}
}

// TestCLI_RootFlagsPassthrough tests that root command accepts visualize flags
// since it runs visualize by default.
func TestCLI_RootFlagsPassthrough(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "only-network flag on root with help",
			args: []string{"--only-network", "bridge", "--help"},
		},
		{
			name: "container flag on root with help",
			args: []string{"--container", "test", "--help"},
		},
		{
			name: "no-aliases flag on root with help",
			args: []string{"--no-aliases", "--help"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"run", ".."}, tc.args...)
			cmd := exec.Command("go", args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("root command with flags %v should be accepted: %v\nstderr: %s", tc.args, err, stderr.String())
			}
		})
	}
}
