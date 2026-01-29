// Package output provides formatters for Docker network visualization.
// This file contains color support utilities.
package output

import (
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var (
	// Color definitions for output formatting
	colorNetwork   = color.New(color.FgCyan, color.Bold)
	colorContainer = color.New(color.FgGreen)
	colorAlias     = color.New(color.FgYellow)
	colorLabel     = color.New(color.FgMagenta)
	colorTree      = color.New(color.FgBlue)
)

// ColorWriter wraps an io.Writer and provides colored output methods.
type ColorWriter struct {
	writer  io.Writer
	enabled bool
}

// NewColorWriter creates a new ColorWriter.
// It detects if color should be enabled based on terminal support and configuration.
func NewColorWriter(w io.Writer) *ColorWriter {
	enabled := shouldUseColor(w)

	// Respect user configuration
	if viper.GetBool("no-color") {
		enabled = false
	}

	return &ColorWriter{
		writer:  w,
		enabled: enabled,
	}
}

// shouldUseColor determines if colored output should be used.
func shouldUseColor(w io.Writer) bool {
	// Check if output is a terminal
	if f, ok := w.(*os.File); ok {
		return !color.NoColor && isTerminal(f)
	}
	return false
}

// isTerminal checks if the file descriptor is a terminal.
func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Network prints text in network color (cyan, bold).
func (cw *ColorWriter) Network(text string) string {
	if !cw.enabled {
		return text
	}
	return colorNetwork.Sprint(text)
}

// Container prints text in container color (green).
func (cw *ColorWriter) Container(text string) string {
	if !cw.enabled {
		return text
	}
	return colorContainer.Sprint(text)
}

// Alias prints text in alias color (yellow).
func (cw *ColorWriter) Alias(text string) string {
	if !cw.enabled {
		return text
	}
	return colorAlias.Sprint(text)
}

// Label prints text in label color (magenta).
func (cw *ColorWriter) Label(text string) string {
	if !cw.enabled {
		return text
	}
	return colorLabel.Sprint(text)
}

// Tree prints text in tree color (blue).
func (cw *ColorWriter) Tree(text string) string {
	if !cw.enabled {
		return text
	}
	return colorTree.Sprint(text)
}

// IsEnabled returns whether color is enabled.
func (cw *ColorWriter) IsEnabled() bool {
	return cw.enabled
}

// Writer returns the underlying writer.
func (cw *ColorWriter) Writer() io.Writer {
	return cw.writer
}
