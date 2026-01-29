package output

import (
	"bytes"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func TestNewColorWriter_WithBuffer(t *testing.T) {
	var buf bytes.Buffer
	cw := NewColorWriter(&buf)

	// Buffer is not a terminal, so color should be disabled
	if cw.enabled {
		t.Error("expected color to be disabled for buffer writer")
	}

	if cw.writer != &buf {
		t.Error("expected writer to be set correctly")
	}
}

func TestNewColorWriter_RespectsNoColorConfig(t *testing.T) {
	// Save original value
	original := viper.GetBool("no-color")
	defer viper.Set("no-color", original)

	// Set no-color to true
	viper.Set("no-color", true)

	var buf bytes.Buffer
	cw := NewColorWriter(&buf)

	if cw.enabled {
		t.Error("expected color to be disabled when no-color config is true")
	}
}

func TestColorWriter_Network_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	result := cw.Network("test-network")

	if result != "test-network" {
		t.Errorf("expected plain text 'test-network', got %q", result)
	}
}

func TestColorWriter_Network_Enabled(t *testing.T) {
	// Force color output for testing
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	result := cw.Network("test-network")

	// When enabled, result should contain ANSI codes
	if result == "test-network" {
		t.Error("expected colored output when enabled")
	}

	// Result should still contain the text
	if len(result) <= len("test-network") {
		t.Error("expected result to contain color codes")
	}
}

func TestColorWriter_Container_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	result := cw.Container("my-container")

	if result != "my-container" {
		t.Errorf("expected plain text 'my-container', got %q", result)
	}
}

func TestColorWriter_Container_Enabled(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	result := cw.Container("my-container")

	if result == "my-container" {
		t.Error("expected colored output when enabled")
	}
}

func TestColorWriter_Alias_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	result := cw.Alias("my-alias")

	if result != "my-alias" {
		t.Errorf("expected plain text 'my-alias', got %q", result)
	}
}

func TestColorWriter_Alias_Enabled(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	result := cw.Alias("my-alias")

	if result == "my-alias" {
		t.Error("expected colored output when enabled")
	}
}

func TestColorWriter_Label_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	result := cw.Label("Network:")

	if result != "Network:" {
		t.Errorf("expected plain text 'Network:', got %q", result)
	}
}

func TestColorWriter_Label_Enabled(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	result := cw.Label("Network:")

	if result == "Network:" {
		t.Error("expected colored output when enabled")
	}
}

func TestColorWriter_Tree_Disabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	result := cw.Tree("├──")

	if result != "├──" {
		t.Errorf("expected plain text '├──', got %q", result)
	}
}

func TestColorWriter_Tree_Enabled(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	result := cw.Tree("├──")

	if result == "├──" {
		t.Error("expected colored output when enabled")
	}
}

func TestColorWriter_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cw := &ColorWriter{writer: &buf, enabled: tt.enabled}

			if cw.IsEnabled() != tt.expected {
				t.Errorf("expected IsEnabled() = %v, got %v", tt.expected, cw.IsEnabled())
			}
		})
	}
}

func TestColorWriter_Writer(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	if cw.Writer() != &buf {
		t.Error("expected Writer() to return the underlying writer")
	}
}

func TestShouldUseColor_NonFileWriter(t *testing.T) {
	var buf bytes.Buffer

	result := shouldUseColor(&buf)

	if result {
		t.Error("expected shouldUseColor to return false for non-file writer")
	}
}

func TestShouldUseColor_WithFile(t *testing.T) {
	// Create a temporary file (not a terminal)
	tmpFile, err := os.CreateTemp("", "color_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Temp file is not a terminal, so should return false
	result := shouldUseColor(tmpFile)

	// We expect false because a regular file is not a terminal
	if result {
		t.Error("expected shouldUseColor to return false for regular file")
	}
}

func TestIsTerminal_WithRegularFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "terminal_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	result := isTerminal(tmpFile)

	if result {
		t.Error("expected isTerminal to return false for regular file")
	}
}

func TestColorWriter_EmptyString(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	// Test all methods with empty string
	if cw.Network("") != "" {
		t.Error("Network should return empty string when given empty string")
	}
	if cw.Container("") != "" {
		t.Error("Container should return empty string when given empty string")
	}
	if cw.Alias("") != "" {
		t.Error("Alias should return empty string when given empty string")
	}
	if cw.Label("") != "" {
		t.Error("Label should return empty string when given empty string")
	}
	if cw.Tree("") != "" {
		t.Error("Tree should return empty string when given empty string")
	}
}

func TestColorWriter_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	tests := []struct {
		name  string
		input string
	}{
		{"unicode tree", "├──"},
		{"unicode end", "└──"},
		{"unicode vertical", "│"},
		{"mixed", "├── hello └──"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cw.Tree(tt.input) != tt.input {
				t.Errorf("expected %q, got %q", tt.input, cw.Tree(tt.input))
			}
		})
	}
}

func TestNewColorWriter_NilNoColorConfig(t *testing.T) {
	// Reset viper to ensure no-color is not set
	viper.Reset()

	var buf bytes.Buffer
	cw := NewColorWriter(&buf)

	// Should still work even when config is not set
	if cw == nil {
		t.Error("expected non-nil ColorWriter")
	}
}

func TestColorWriter_AllColorMethodsReturnTextWhenDisabled(t *testing.T) {
	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: false}

	testText := "test-text-123"

	methods := []struct {
		name   string
		method func(string) string
	}{
		{"Network", cw.Network},
		{"Container", cw.Container},
		{"Alias", cw.Alias},
		{"Label", cw.Label},
		{"Tree", cw.Tree},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.method(testText)
			if result != testText {
				t.Errorf("%s: expected %q, got %q", m.name, testText, result)
			}
		})
	}
}

func TestColorWriter_AllColorMethodsReturnColoredTextWhenEnabled(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	var buf bytes.Buffer
	cw := &ColorWriter{writer: &buf, enabled: true}

	testText := "test-text-123"

	methods := []struct {
		name   string
		method func(string) string
	}{
		{"Network", cw.Network},
		{"Container", cw.Container},
		{"Alias", cw.Alias},
		{"Label", cw.Label},
		{"Tree", cw.Tree},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.method(testText)
			// When enabled, the result should be longer due to ANSI codes
			if len(result) <= len(testText) {
				t.Errorf("%s: expected colored output longer than input", m.name)
			}
		})
	}
}
