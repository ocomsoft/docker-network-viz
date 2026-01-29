package output

import (
	"testing"
)

func TestTreeSymbolConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "TreeBranch",
			constant: TreeBranch,
			expected: "\u251c\u2500\u2500", // ├──
		},
		{
			name:     "TreeEnd",
			constant: TreeEnd,
			expected: "\u2514\u2500\u2500", // └──
		},
		{
			name:     "TreeVertical",
			constant: TreeVertical,
			expected: "\u2502   ", // │ followed by 3 spaces
		},
		{
			name:     "TreeSpace",
			constant: TreeSpace,
			expected: "    ", // 4 spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.constant)
			}
		})
	}
}

func TestTreeBranchAndEndHaveSameByteLength(t *testing.T) {
	// TreeBranch and TreeEnd should have same byte length for alignment
	if len(TreeBranch) != len(TreeEnd) {
		t.Errorf("TreeBranch (%d bytes) and TreeEnd (%d bytes) should have same length",
			len(TreeBranch), len(TreeEnd))
	}
}

func TestTreeIndentHaveSameByteLength(t *testing.T) {
	// TreeVertical and TreeSpace should have same byte length for alignment
	// Note: TreeVertical contains a Unicode character (│) which is 3 bytes + 3 spaces = 6 bytes
	// TreeSpace is 4 ASCII spaces = 4 bytes
	// However, they have the same VISUAL width (4 characters)
	// This test verifies visual alignment consistency
	verticalRunes := []rune(TreeVertical)
	spaceRunes := []rune(TreeSpace)

	if len(verticalRunes) != len(spaceRunes) {
		t.Errorf("TreeVertical (%d runes) and TreeSpace (%d runes) should have same rune count",
			len(verticalRunes), len(spaceRunes))
	}
}

func TestTreeSymbolsAreUnicode(t *testing.T) {
	// Verify the box-drawing characters are proper Unicode
	branchRunes := []rune(TreeBranch)
	if branchRunes[0] != '\u251c' {
		t.Errorf("TreeBranch first rune should be U+251C, got U+%04X", branchRunes[0])
	}

	endRunes := []rune(TreeEnd)
	if endRunes[0] != '\u2514' {
		t.Errorf("TreeEnd first rune should be U+2514, got U+%04X", endRunes[0])
	}

	verticalRunes := []rune(TreeVertical)
	if verticalRunes[0] != '\u2502' {
		t.Errorf("TreeVertical first rune should be U+2502, got U+%04X", verticalRunes[0])
	}
}

func TestTreeBranchContainsHorizontalLines(t *testing.T) {
	// TreeBranch should contain horizontal line characters
	branchRunes := []rune(TreeBranch)
	if len(branchRunes) < 3 {
		t.Fatalf("TreeBranch should have at least 3 runes, got %d", len(branchRunes))
	}

	// Second and third runes should be horizontal lines (─)
	if branchRunes[1] != '\u2500' || branchRunes[2] != '\u2500' {
		t.Errorf("TreeBranch should have horizontal lines (U+2500), got U+%04X and U+%04X",
			branchRunes[1], branchRunes[2])
	}
}

func TestTreeEndContainsHorizontalLines(t *testing.T) {
	// TreeEnd should contain horizontal line characters
	endRunes := []rune(TreeEnd)
	if len(endRunes) < 3 {
		t.Fatalf("TreeEnd should have at least 3 runes, got %d", len(endRunes))
	}

	// Second and third runes should be horizontal lines (─)
	if endRunes[1] != '\u2500' || endRunes[2] != '\u2500' {
		t.Errorf("TreeEnd should have horizontal lines (U+2500), got U+%04X and U+%04X",
			endRunes[1], endRunes[2])
	}
}

func TestTreeSpaceIsAllSpaces(t *testing.T) {
	for i, r := range TreeSpace {
		if r != ' ' {
			t.Errorf("TreeSpace should contain only spaces, found %q at position %d", r, i)
		}
	}
}
