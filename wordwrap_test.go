package wordwrap

import (
	"reflect"
	"testing"
)

func TestWordWrap(t *testing.T) {
	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
	maxWidth := 40
	expected := []string{
		"Lorem ipsum dolor sit amet, consectetur",
		"adipiscing elit. Sed do eiusmod tempor",
		"incididunt ut labore et dolore magna",
		"aliqua.",
	}

	result, err := WordWrap(text, maxWidth)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestWordWrapInvalidWidth(t *testing.T) {
	_, err := WordWrap("hello", 0)
	if err == nil {
		t.Error("Expected error for maxWidth=0")
	}
	_, err = WordWrap("hello", -1)
	if err == nil {
		t.Error("Expected error for maxWidth=-1")
	}
}

func TestWordWrapEmptyInput(t *testing.T) {
	result, err := WordWrap("", 40)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, []string{""}) {
		t.Errorf("Expected [\"\"], got %v", result)
	}
}

func TestWordWrapMultiLine(t *testing.T) {
	text := "short\nalso short"
	result, err := WordWrap(text, 40)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := []string{"short", "also short"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestWordWrapHardBreak(t *testing.T) {
	// A single long word with no spaces must be hard-broken
	text := "abcdefghijklmnopqrstuvwxyz"
	result, err := WordWrap(text, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := []string{"abcdefghij", "klmnopqrst", "uvwxyz"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// --- NoLineStart tests ---

func TestNoLineStart(t *testing.T) {
	// Characters that must NOT start a line
	for _, r := range []rune{'.', ',', '!', '?', ':', ';', ')', ']', '}', '\'', '"',
		'\u2019', '\u201D', '\u00BB', '\u2014', '\u2013'} {
		if !NoLineStart(r) {
			t.Errorf("NoLineStart(%q) should be true", r)
		}
	}
	// Characters that CAN start a line
	for _, r := range []rune{'a', 'Z', '0', '(', '[', '{', ' ', '-', '_'} {
		if NoLineStart(r) {
			t.Errorf("NoLineStart(%q) should be false", r)
		}
	}
}

// --- WrapLine tests ---

func TestWrapLineWithinLimit(t *testing.T) {
	r := WrapLine("hello world", 20, 0)
	if r.Wrapped {
		t.Error("Should not wrap a line within the limit")
	}
	if r.Left != "hello world" {
		t.Errorf("Left should be the full line, got %q", r.Left)
	}
	if r.BreakAt != -1 {
		t.Errorf("BreakAt should be -1, got %d", r.BreakAt)
	}
}

func TestWrapLineExactlyAtLimit(t *testing.T) {
	// 10 chars, limit 10 — should NOT wrap
	r := WrapLine("1234567890", 10, 0)
	if r.Wrapped {
		t.Error("Should not wrap a line exactly at the limit")
	}
}

func TestWrapLineOneOverLimit(t *testing.T) {
	// "hello world" = 11 chars, limit 10
	r := WrapLine("hello world", 10, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap a line one over the limit")
	}
	if r.Left != "hello" {
		t.Errorf("Left = %q, want %q", r.Left, "hello")
	}
	if r.Right != "world" {
		t.Errorf("Right = %q, want %q", r.Right, "world")
	}
	if r.BreakAt != 5 {
		t.Errorf("BreakAt = %d, want 5", r.BreakAt)
	}
}

func TestWrapLineSpaceAtExactLimit(t *testing.T) {
	// Space at exactly position 10 (the limit). Left should be 10 chars.
	// "abcdefghij klm" — space at index 10
	line := "abcdefghij klm"
	r := WrapLine(line, 10, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "abcdefghij" {
		t.Errorf("Left = %q, want %q", r.Left, "abcdefghij")
	}
	if r.Right != "klm" {
		t.Errorf("Right = %q, want %q", r.Right, "klm")
	}
	if r.BreakAt != 10 {
		t.Errorf("BreakAt = %d, want 10", r.BreakAt)
	}
}

func TestWrapLinePreferLatestSpace(t *testing.T) {
	// Two spaces: at 3 and 7. Limit 8. Should break at 7 (latest valid space).
	// "abc def ghij" — spaces at 3 and 7
	line := "abc def ghij"
	r := WrapLine(line, 8, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "abc def" {
		t.Errorf("Left = %q, want %q", r.Left, "abc def")
	}
	if r.Right != "ghij" {
		t.Errorf("Right = %q, want %q", r.Right, "ghij")
	}
}

func TestWrapLineNoSpaces(t *testing.T) {
	// No spaces at all — can't wrap
	r := WrapLine("abcdefghijklmno", 10, 0)
	if r.Wrapped {
		t.Error("Should not wrap when there are no spaces")
	}
}

func TestWrapLineTypographyPeriod(t *testing.T) {
	// Don't break before a period: "word .next" — space at 4, but next char is '.'
	// Should NOT break there.
	line := "word .next more"
	r := WrapLine(line, 6, 0)
	// The only space within range 1..6 is at index 4, but '.' follows it.
	// So no valid break point.
	if r.Wrapped {
		t.Error("Should not wrap before a period")
	}
}

func TestWrapLineTypographyComma(t *testing.T) {
	// "word ,next more" — don't break before comma
	line := "word ,next more"
	r := WrapLine(line, 6, 0)
	if r.Wrapped {
		t.Error("Should not wrap before a comma")
	}
}

func TestWrapLineTypographySkipsToEarlierBreak(t *testing.T) {
	// "ab cd .ef gh" — space at 2, 5, 9. Limit 7.
	// Space at 5 followed by '.', so skip it. Fall back to space at 2.
	line := "ab cd .ef gh"
	r := WrapLine(line, 7, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "ab" {
		t.Errorf("Left = %q, want %q", r.Left, "ab")
	}
	// Right should be "cd .ef gh" — the period stays grouped with its preceding word
	if r.Right != "cd .ef gh" {
		t.Errorf("Right = %q, want %q", r.Right, "cd .ef gh")
	}
}

func TestWrapLineTypographyClosingParen(t *testing.T) {
	// "call (foo )bar end" — space at 4, 9. Position 10 is ')'.
	// Breaking at 9 would start the new line with ')'. Should skip.
	// Falls back to space at 4.
	line := "call (foo )bar end"
	r := WrapLine(line, 11, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "call" {
		t.Errorf("Left = %q, want %q", r.Left, "call")
	}
}

func TestWrapLineConsecutiveSpaces(t *testing.T) {
	// "hello   world!" — spaces at 5, 6, 7. Should break cleanly.
	line := "hello   world!"
	r := WrapLine(line, 9, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	// Break at position 7 (last space in run, closest to limit).
	// Left part has trailing spaces trimmed.
	if r.Left != "hello" {
		t.Errorf("Left = %q, want %q", r.Left, "hello")
	}
	if r.Right != "world!" {
		t.Errorf("Right = %q, want %q", r.Right, "world!")
	}
}

func TestWrapLineTabAsWhitespace(t *testing.T) {
	// Tab character should be recognized as a valid break point
	line := "hello\tworld and more"
	r := WrapLine(line, 8, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap at tab")
	}
	if r.Left != "hello" {
		t.Errorf("Left = %q, want %q", r.Left, "hello")
	}
	if r.Right != "world and more" {
		t.Errorf("Right = %q, want %q", r.Right, "world and more")
	}
}

func TestWrapLineUnicode(t *testing.T) {
	// Unicode characters: each emoji is one rune
	line := "héllo wörld thïs ïs ünïcödé"
	r := WrapLine(line, 15, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "héllo wörld" {
		t.Errorf("Left = %q, want %q", r.Left, "héllo wörld")
	}
	if r.Right != "thïs ïs ünïcödé" {
		t.Errorf("Right = %q, want %q", r.Right, "thïs ïs ünïcödé")
	}
}

func TestWrapLineMaxBacktrack(t *testing.T) {
	// "aa bb cc dd ee ff" — spaces at 2, 5, 8, 11, 14. Limit 12.
	// With maxBacktrack 4, search from 12 down to 8.
	// Space at 11 is within range. Break there.
	line := "aa bb cc dd ee ff"
	r := WrapLine(line, 12, 4)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "aa bb cc dd" {
		t.Errorf("Left = %q, want %q", r.Left, "aa bb cc dd")
	}
}

func TestWrapLineMaxBacktrackTooRestrictive(t *testing.T) {
	// "hello world" — space at 5. Limit 10, maxBacktrack 2.
	// Search from 10 down to 8. No space in range 8..10.
	line := "hello worldx"
	r := WrapLine(line, 10, 2)
	if r.Wrapped {
		t.Error("Should not wrap when maxBacktrack is too restrictive")
	}
}

func TestWrapLineOnlySpaceAtPosition0(t *testing.T) {
	// Space only at position 0 — should not break (would produce empty left)
	line := " abcdefghij"
	r := WrapLine(line, 5, 0)
	if r.Wrapped {
		t.Error("Should not break at position 0 (empty left part)")
	}
}

func TestWrapLineZeroLimit(t *testing.T) {
	r := WrapLine("hello world", 0, 0)
	if r.Wrapped {
		t.Error("Should not wrap with limit 0")
	}
}

func TestWrapLineNegativeLimit(t *testing.T) {
	r := WrapLine("hello world", -5, 0)
	if r.Wrapped {
		t.Error("Should not wrap with negative limit")
	}
}

// --- Wrap-as-you-type scenario tests ---

func TestWrapAsYouTypeGitCommit(t *testing.T) {
	// Simulates typing a git commit message with limit 72.
	// The line just went over 72 characters.
	line := "Fix the rendering issue that caused the sidebar to overlap with the main content area"
	r := WrapLine(line, 72, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	leftRunes := []rune(r.Left)
	if len(leftRunes) > 72 {
		t.Errorf("Left part has %d runes, should be <= 72", len(leftRunes))
	}
	if len(r.Right) == 0 {
		t.Error("Right part should not be empty")
	}
	// Verify the content is preserved
	if r.Left+" "+r.Right != line {
		t.Errorf("Content not preserved: %q + %q != %q", r.Left, r.Right, line)
	}
}

func TestWrapAsYouTypeEmail(t *testing.T) {
	// Simulates typing an email with limit 72
	line := "Thank you for your email regarding the project timeline and deliverables schedule"
	r := WrapLine(line, 72, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	leftRunes := []rune(r.Left)
	if len(leftRunes) > 72 {
		t.Errorf("Left part has %d runes, should be <= 72", len(leftRunes))
	}
}

func TestWrapAsYouTypeExactLimit(t *testing.T) {
	// Line is exactly at the limit — should NOT wrap
	// 72 characters exactly:
	line := "This line has exactly seventy-two characters and should not be wrapped!"
	runes := []rune(line)
	limit := len(runes)
	r := WrapLine(line, limit, 0)
	if r.Wrapped {
		t.Error("Should not wrap a line that is exactly at the limit")
	}
}

func TestWrapAsYouTypeOneCharOver(t *testing.T) {
	// Line is exactly one character over the limit — should wrap
	line := "This is a test line that goes just one character over the limit right x"
	limit := len([]rune(line)) - 1
	r := WrapLine(line, limit, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap when one char over limit")
	}
	if len([]rune(r.Left)) > limit {
		t.Errorf("Left part exceeds limit: %d > %d", len([]rune(r.Left)), limit)
	}
}

func TestWrapAsYouTypeSpaceAtLimit(t *testing.T) {
	// When a space falls exactly at the limit position, it should be used
	// as the break point, producing a left part of exactly `limit` chars.
	// Build a line: 10 non-space chars + space at position 10 + more chars
	line := "0123456789 abcde"
	r := WrapLine(line, 10, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	if r.Left != "0123456789" {
		t.Errorf("Left = %q, want %q", r.Left, "0123456789")
	}
	if r.BreakAt != 10 {
		t.Errorf("BreakAt = %d, want 10", r.BreakAt)
	}
	if r.Right != "abcde" {
		t.Errorf("Right = %q, want %q", r.Right, "abcde")
	}
}

func TestWrapAsYouTypeTypographyPreserved(t *testing.T) {
	// When wrapping, the new line must not start with a period.
	// "some text here. More text follows" with limit that would break before '.'
	line := "text here .More stuff after"
	r := WrapLine(line, 12, 0)
	if r.Wrapped {
		rightRunes := []rune(r.Right)
		if len(rightRunes) > 0 && NoLineStart(rightRunes[0]) {
			t.Errorf("Right part starts with no-line-start char %q", rightRunes[0])
		}
	}
}

func TestWrapAsYouTypePreservesContent(t *testing.T) {
	// Verify that no content is lost during wrapping
	line := "The quick brown fox jumps over the lazy dog near the riverbank"
	r := WrapLine(line, 30, 0)
	if !r.Wrapped {
		t.Fatal("Should wrap")
	}
	// Reconstruct: Left + " " + Right should equal original
	// (assuming single space at break)
	reconstructed := r.Left + " " + r.Right
	if reconstructed != line {
		t.Errorf("Content not preserved:\noriginal:      %q\nreconstructed: %q", line, reconstructed)
	}
}

func TestWrapLineEmDash(t *testing.T) {
	// Em dash should not start a new line
	line := "word —continuation more text here"
	r := WrapLine(line, 6, 0)
	// Space at 4, but next is '—' (noLineStart). Should not break there.
	if r.Wrapped {
		rightRunes := []rune(r.Right)
		if len(rightRunes) > 0 && NoLineStart(rightRunes[0]) {
			t.Errorf("Right starts with em dash")
		}
	}
}

func TestWrapLineRightGuillemet(t *testing.T) {
	// Right guillemet should not start a new line
	line := "word »quote more text here"
	r := WrapLine(line, 6, 0)
	if r.Wrapped {
		rightRunes := []rune(r.Right)
		if len(rightRunes) > 0 && NoLineStart(rightRunes[0]) {
			t.Errorf("Right starts with right guillemet")
		}
	}
}

func TestWordWrapPreservesNewlines(t *testing.T) {
	text := "short\n\nanother paragraph that is long enough to need wrapping at the limit"
	result, err := WordWrap(text, 30)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Should have at least 4 lines: "short", "", wrapped parts of paragraph
	if len(result) < 4 {
		t.Errorf("Expected at least 4 lines, got %d: %v", len(result), result)
	}
	// Empty line should be preserved
	if result[1] != "" {
		t.Errorf("Expected empty line preserved, got %q", result[1])
	}
}

func TestWordWrapTypographyNotBrokenBeforePeriod(t *testing.T) {
	text := "This is a sentence .that should not break before the period character"
	result, err := WordWrap(text, 22)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	for _, line := range result {
		runes := []rune(line)
		if len(runes) > 0 && NoLineStart(runes[0]) {
			t.Errorf("Line starts with no-line-start char: %q", line)
		}
	}
}
