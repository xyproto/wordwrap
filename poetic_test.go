package wordwrap

import (
	"reflect"
	"strings"
	"testing"
)

func TestPoeticWrapInvalidWidth(t *testing.T) {
	if _, err := PoeticWrap("hello", 0, 0); err == nil {
		t.Error("Expected error for maxWidth=0")
	}
	if _, err := PoeticWrap("hello", -1, 0); err == nil {
		t.Error("Expected error for maxWidth=-1")
	}
}

func TestPoeticWrapEmpty(t *testing.T) {
	result, err := PoeticWrap("", 40, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, []string{""}) {
		t.Errorf("Expected [\"\"], got %v", result)
	}
}

func TestPoeticWrapShortFits(t *testing.T) {
	// Fits on one line.
	result, err := PoeticWrap("Be brave today.", 40, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, []string{"Be brave today."}) {
		t.Errorf("Got %v", result)
	}
}

func TestPoeticWrapPrefersComma(t *testing.T) {
	// Motivating fortunecraft case: break at the comma.
	text := "The future of humanity lies in the belly of a giant space potato, " +
		"which will explode with the force of a thousand exploding potatoes!"
	expected := []string{
		"The future of humanity lies in the belly of a giant space potato,",
		"which will explode with the force of a thousand exploding potatoes!",
	}
	result, err := PoeticWrap(text, 80, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q\n     got %q", expected, result)
	}
}

func TestPoeticWrapPrefersCommaAtNarrowerWidth(t *testing.T) {
	// Comma is still in range at the narrower width.
	text := "The future of humanity lies in the belly of a giant space potato, " +
		"which will explode with the force of a thousand exploding potatoes!"
	result, err := PoeticWrap(text, 72, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) < 2 {
		t.Fatalf("Expected at least 2 lines, got %v", result)
	}
	if !strings.HasSuffix(result[0], "potato,") {
		t.Errorf("Expected first line to end at the comma, got %q", result[0])
	}
}

func TestPoeticWrapPrefersSentenceEnd(t *testing.T) {
	// Break at the period rather than mid-clause.
	text := "Wisdom comes to those who listen carefully. Speak softly and often."
	result, err := PoeticWrap(text, 50, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) < 2 {
		t.Fatalf("Expected wrapping, got %v", result)
	}
	if !strings.HasSuffix(result[0], "carefully.") {
		t.Errorf("Expected first line to end at the period, got %q", result[0])
	}
}

func TestPoeticWrapNoNaturalBreak(t *testing.T) {
	// No punctuation: fall back to greedy word boundaries.
	text := "the quick brown fox jumps over the lazy dog and runs through fields"
	result, err := PoeticWrap(text, 30, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) < 2 {
		t.Fatalf("Expected wrapping, got %v", result)
	}
	for _, line := range result {
		if len([]rune(line)) > 30 {
			t.Errorf("Line %q exceeds maxWidth", line)
		}
	}
}

func TestPoeticWrapRespectsMaxWidth(t *testing.T) {
	// No line of normal-sized words should exceed maxWidth.
	text := "Alpha, beta, gamma, delta; epsilon, zeta, eta, theta. Iota, kappa!"
	const maxWidth = 25
	result, err := PoeticWrap(text, maxWidth, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	for _, line := range result {
		if len([]rune(line)) > maxWidth {
			t.Errorf("Line %q exceeds maxWidth=%d", line, maxWidth)
		}
	}
}

func TestPoeticWrapTypographyNoLineStart(t *testing.T) {
	// No line may start with a NoLineStart char.
	text := "first part of the line .next bit comes after with more words still"
	result, err := PoeticWrap(text, 20, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	for _, line := range result {
		runes := []rune(line)
		if len(runes) > 0 && NoLineStart(runes[0]) {
			t.Errorf("Line starts with NoLineStart char: %q", line)
		}
	}
}

func TestPoeticWrapPreservesNewlines(t *testing.T) {
	// Newlines are hard paragraph breaks; blanks preserved.
	text := "Short stanza one.\n\nA longer second stanza that will need to wrap across two lines."
	result, err := PoeticWrap(text, 40, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) < 4 {
		t.Fatalf("Expected at least 4 lines, got %v", result)
	}
	if result[0] != "Short stanza one." {
		t.Errorf("Expected first stanza preserved, got %q", result[0])
	}
	if result[1] != "" {
		t.Errorf("Expected blank line between stanzas, got %q", result[1])
	}
}

func TestPoeticWrapSingleOversizedWord(t *testing.T) {
	// Oversized word survives on its own line, not hard-broken.
	text := "supercalifragilisticexpialidocious is fun"
	result, err := PoeticWrap(text, 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) == 0 || result[0] != "supercalifragilisticexpialidocious" {
		t.Errorf("Oversized word not preserved, got %v", result)
	}
}

func TestPoeticWrapMinWidthDefault(t *testing.T) {
	// Default minWidth (maxWidth*2/5) blocks the comma at column 3.
	text := "Hi, " + strings.Repeat("word ", 30)
	result, err := PoeticWrap(text, 50, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if strings.HasSuffix(result[0], "Hi,") {
		t.Errorf("Comma at column 3 should be below minWidth, got %q", result[0])
	}
}

func TestPoeticWrapMinWidthAboveMax(t *testing.T) {
	// minWidth > maxWidth disables the guard; must still complete.
	text := "Alpha, beta, gamma, delta. Epsilon, zeta, eta, theta!"
	const maxWidth = 30
	result, err := PoeticWrap(text, maxWidth, 1000)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("Expected at least one line")
	}
	for _, line := range result {
		if len([]rune(line)) > maxWidth {
			t.Errorf("Line %q exceeds maxWidth=%d", line, maxWidth)
		}
	}
}

func TestPoeticWrapPreservesContent(t *testing.T) {
	// Rejoined lines must reproduce the original words in order.
	text := "Alpha, beta, gamma; delta — epsilon: zeta. Eta theta iota kappa!"
	result, err := PoeticWrap(text, 20, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got := strings.Join(result, " ")
	want := strings.Join(strings.Fields(text), " ")
	if got != want {
		t.Errorf("Content not preserved:\nwant %q\n got %q", want, got)
	}
}

func TestPoeticWrapUnicode(t *testing.T) {
	// Multi-byte runes count as one column each.
	text := "Café au lait, étoile filante; rêverie éphémère."
	const maxWidth = 25
	result, err := PoeticWrap(text, maxWidth, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	for _, line := range result {
		if len([]rune(line)) > maxWidth {
			t.Errorf("Line %q exceeds maxWidth=%d (runes=%d)", line, maxWidth, len([]rune(line)))
		}
	}
}

func TestPoeticWrapEmDash(t *testing.T) {
	// Em dash near the right margin should be chosen.
	text := "After many quiet years she finally returned — to find the village gone."
	result, err := PoeticWrap(text, 50, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) < 2 {
		t.Fatalf("Expected wrapping, got %v", result)
	}
	if !strings.HasSuffix(result[0], "—") {
		t.Errorf("Expected first line to end at the em dash, got %q", result[0])
	}
}

func TestBreakPenaltyOrdering(t *testing.T) {
	// Penalty ordering drives the DP's poetic preferences.
	cases := []struct {
		word   string
		expect float64
	}{
		{"end.", breakSentenceEnd},
		{"end!", breakSentenceEnd},
		{"end?", breakSentenceEnd},
		{"end;", breakSemicolon},
		{"end:", breakColon},
		{"end,", breakComma},
		{"end—", breakDash},
		{"end", breakMidClause},
		{"(end.)", breakSentenceEnd},
		{"\"end!\"", breakSentenceEnd},
	}
	for _, c := range cases {
		if got := breakPenalty(c.word); got != c.expect {
			t.Errorf("breakPenalty(%q) = %v, want %v", c.word, got, c.expect)
		}
	}
	if !(breakSentenceEnd < breakSemicolon &&
		breakSemicolon < breakColon &&
		breakColon < breakComma &&
		breakComma < breakDash &&
		breakDash < breakMidClause) {
		t.Error("Break penalty constants are not strictly ordered")
	}
}
