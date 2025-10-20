package tui

import (
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestTruncateTodoTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		maxLength int
		expected  string
	}{
		{
			name:      "Short title - no truncation",
			title:     "Short",
			maxLength: 10,
			expected:  "Short",
		},
		{
			name:      "Exactly max length - no truncation",
			title:     "1234567890",
			maxLength: 10,
			expected:  "1234567890",
		},
		{
			name:      "Long title - truncate with ellipsis",
			title:     "This is a very long todo title",
			maxLength: 15,
			expected:  "This is a very ...",
		},
		{
			name:      "Empty title",
			title:     "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "Japanese characters - no truncation",
			title:     "買い物に行く",
			maxLength: 10,
			expected:  "買い物に行く",
		},
		{
			name:      "Japanese characters - truncation",
			title:     "これは非常に長いToDoのタイトルです",
			maxLength: 15,
			expected:  "これは非常に長いToDoのタイ...",
		},
		{
			name:      "Mixed English and Japanese - truncation",
			title:     "Buy groceries at スーパーマーケット",
			maxLength: 15,
			expected:  "Buy groceries a...",
		},
		{
			name:      "Emoji characters",
			title:     "🎉🎊🎈🎁🎂🍰🍕🍔🍟🌮",
			maxLength: 5,
			expected:  "🎉🎊🎈🎁🎂...",
		},
		{
			name:      "Max length 1",
			title:     "ABC",
			maxLength: 1,
			expected:  "A...",
		},
		{
			name:      "Max length 0 with content",
			title:     "Test",
			maxLength: 0,
			expected:  "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateTodoTitle(tt.title, tt.maxLength)
			if result != tt.expected {
				t.Errorf("truncateTodoTitle(%q, %d) = %q; expected %q",
					tt.title, tt.maxLength, result, tt.expected)
			}

			// Additional check: result should not be longer than maxLength + 3 (for "...")
			runes := []rune(result)
			if len(runes) > tt.maxLength+3 {
				t.Errorf("truncateTodoTitle(%q, %d) returned %q with length %d; should not exceed %d",
					tt.title, tt.maxLength, result, len(runes), tt.maxLength+3)
			}
		})
	}
}

func TestTruncateTodoTitle_RuneCount(t *testing.T) {
	// This test specifically verifies that we count runes, not bytes
	tests := []struct {
		name      string
		title     string
		maxLength int
		wantRunes int // Expected number of runes in result (excluding "...")
	}{
		{
			name:      "ASCII only",
			title:     "Hello World",
			maxLength: 5,
			wantRunes: 5, // "Hello" + "..."
		},
		{
			name:      "Japanese only",
			title:     "こんにちは世界",
			maxLength: 5,
			wantRunes: 5, // "こんにちは" + "..."
		},
		{
			name:      "Emoji only",
			title:     "😀😁😂😃😄😅😆",
			maxLength: 3,
			wantRunes: 3, // "😀😁😂" + "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateTodoTitle(tt.title, tt.maxLength)
			runes := []rune(result)

			// If truncated, result should be maxLength + 3 ("...")
			if len([]rune(tt.title)) > tt.maxLength {
				expectedLength := tt.maxLength + 3
				if len(runes) != expectedLength {
					t.Errorf("truncateTodoTitle(%q, %d) = %q with %d runes; expected %d runes",
						tt.title, tt.maxLength, result, len(runes), expectedLength)
				}
			}
		})
	}
}

func TestTruncateStringByWidth(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxWidth  int
	}{
		{"Short ASCII", "test", 10},
		{"Long ASCII", "This is a very long string", 10},
		{"Short Japanese", "テスト", 10},
		{"Long Japanese", "これは非常に長い文字列です", 10},
		{"Mixed short", "test テスト", 15},
		{"Exact width ASCII", "12345", 5},
		{"Exact width Japanese", "あい", 4},
		{"Empty string", "", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateStringByWidth(tt.input, tt.maxWidth)
			actualWidth := runewidth.StringWidth(result)

			if actualWidth > tt.maxWidth {
				t.Errorf("truncateStringByWidth(%q, %d) = %q (width %d), exceeds maxWidth %d",
					tt.input, tt.maxWidth, result, actualWidth, tt.maxWidth)
			}

			originalWidth := runewidth.StringWidth(tt.input)
			if originalWidth <= tt.maxWidth {
				if result != tt.input {
					t.Errorf("truncateStringByWidth(%q, %d) = %q, want %q (no truncation)",
						tt.input, tt.maxWidth, result, tt.input)
				}
			} else {
				if len(result) < 3 || result[len(result)-3:] != "..." {
					t.Errorf("truncateStringByWidth(%q, %d) = %q, should end with '...'",
						tt.input, tt.maxWidth, result)
				}
			}
		})
	}
}

func TestPadStringToWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
	}{
		{"ASCII short", "test", 10},
		{"ASCII exact", "test", 4},
		{"ASCII long", "testing", 5},
		{"Japanese short", "テスト", 10},
		{"Japanese exact", "あい", 4},
		{"Mixed", "test あ", 10},
		{"Empty", "", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padStringToWidth(tt.input, tt.width)
			actualWidth := runewidth.StringWidth(result)

			inputWidth := runewidth.StringWidth(tt.input)
			expectedWidth := tt.width
			if inputWidth > tt.width {
				expectedWidth = inputWidth
			}

			if actualWidth != expectedWidth {
				t.Errorf("padStringToWidth(%q, %d) width = %d, want %d. Result: %q",
					tt.input, tt.width, actualWidth, expectedWidth, result)
			}

			if actualWidth < inputWidth {
				t.Errorf("padStringToWidth(%q, %d) resulted in shorter width %d < %d",
					tt.input, tt.width, actualWidth, inputWidth)
			}
		})
	}
}
