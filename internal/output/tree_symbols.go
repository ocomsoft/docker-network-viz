// Package output provides formatters for Docker network visualization.
// This file contains tree drawing symbol constants.
package output

// Tree drawing symbols for consistent output formatting
const (
	// TreeBranch is the branch symbol for non-last items
	TreeBranch = "\u251c\u2500\u2500" // ├──

	// TreeEnd is the end symbol for the last item
	TreeEnd = "\u2514\u2500\u2500" // └──

	// TreeVertical is the vertical line symbol for continuing branches
	TreeVertical = "\u2502   " // │ followed by spaces

	// TreeSpace is the indent for items after the last branch
	TreeSpace = "    "
)
