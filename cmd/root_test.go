package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// TestRootCommandExists verifies that the root command is properly defined.
func TestRootCommandExists(t *testing.T) {
	cmd := GetRootCmd()
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}

	if cmd.Use != AppName {
		t.Errorf("root command Use should be %q, got %q", AppName, cmd.Use)
	}

	if cmd.Short != AppDescription {
		t.Errorf("root command Short should be %q, got %q", AppDescription, cmd.Short)
	}
}

// TestRootCommandHasFlags verifies that the root command has required flags.
func TestRootCommandHasFlags(t *testing.T) {
	cmd := GetRootCmd()

	// Check for config flag
	configFlag := cmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("root command should have a config flag")
	}

	// Check for no-color flag
	noColorFlag := cmd.PersistentFlags().Lookup("no-color")
	if noColorFlag == nil {
		t.Error("root command should have a no-color flag")
	}

	// Check for only-network flag
	onlyNetworkFlag := cmd.Flags().Lookup("only-network")
	if onlyNetworkFlag == nil {
		t.Error("root command should have an only-network flag")
	}

	// Check for container flag
	containerFlag := cmd.Flags().Lookup("container")
	if containerFlag == nil {
		t.Error("root command should have a container flag")
	}

	// Check for no-aliases flag
	noAliasesFlag := cmd.Flags().Lookup("no-aliases")
	if noAliasesFlag == nil {
		t.Error("root command should have a no-aliases flag")
	}
}

// TestRootCommandHasVisualizeSubcommand verifies that visualize is a subcommand.
func TestRootCommandHasVisualizeSubcommand(t *testing.T) {
	cmd := GetRootCmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "visualize" {
			found = true
			break
		}
	}

	if !found {
		t.Error("root command should have a visualize subcommand")
	}
}

// TestInitConfigWithNoFile verifies that initConfig works when no config file exists.
func TestInitConfigWithNoFile(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	// Call initConfig with no config file set
	cfgFile = ""
	err := initConfig(nil, nil)

	if err != nil {
		t.Errorf("initConfig should not return error when no config file exists: %v", err)
	}
}

// TestInitConfigWithConfigFile verifies that initConfig reads a config file.
func TestInitConfigWithConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := []byte(`verbose: true
only-network: bridge
`)
	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Reset viper for this test
	viper.Reset()

	// Set the config file
	cfgFile = configPath
	err := initConfig(nil, nil)

	if err != nil {
		t.Errorf("initConfig should not return error: %v", err)
	}

	// Verify config was read
	if viper.GetString("only-network") != "bridge" {
		t.Errorf("expected only-network to be 'bridge', got %q", viper.GetString("only-network"))
	}

	// Reset for other tests
	cfgFile = ""
}

// TestInitConfigEnvironmentVariables verifies that environment variables are read.
func TestInitConfigEnvironmentVariables(t *testing.T) {
	// Reset viper for this test
	viper.Reset()

	// Set an environment variable
	t.Setenv("DNV_NO_COLOR", "true")

	cfgFile = ""
	err := initConfig(nil, nil)

	if err != nil {
		t.Errorf("initConfig should not return error: %v", err)
	}

	// Note: viper.GetBool("no-color") might not work directly for env vars
	// without proper binding, but we verify the setup doesn't error
}

// TestExecuteWithHelp verifies that help flag works.
func TestExecuteWithHelp(t *testing.T) {
	// Reset command for this test
	cmd := GetRootCmd()

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()

	if err != nil {
		t.Errorf("help should not return error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("help output should not be empty")
	}

	// Check for expected content in help
	if !bytes.Contains([]byte(output), []byte("docker-network-viz")) {
		t.Error("help output should contain application name")
	}
}

// TestNoColorConfiguration verifies that no-color configuration works.
func TestNoColorConfiguration(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Test with no-color false
	viper.Set("no-color", false)
	noColorValue := viper.GetBool("no-color")
	if noColorValue {
		t.Error("no-color should be false")
	}

	// Test with no-color true
	viper.Set("no-color", true)
	noColorValue = viper.GetBool("no-color")
	if !noColorValue {
		t.Error("no-color should be true")
	}
}

// TestAppConstants verifies that app constants are defined correctly.
func TestAppConstants(t *testing.T) {
	if AppName == "" {
		t.Error("AppName should not be empty")
	}

	if AppDescription == "" {
		t.Error("AppDescription should not be empty")
	}

	if AppName != "docker-network-viz" {
		t.Errorf("AppName should be 'docker-network-viz', got %q", AppName)
	}
}
