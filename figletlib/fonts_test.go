package figletlib

import (
	"strings"
	"testing"
)

func TestReadFontFromBytesAcceptsNegativeCodeTags(t *testing.T) {
	font, err := ReadFontFromBytes([]byte(minimalFontWithCodeTag("-5", "T")))
	if err != nil {
		t.Fatalf("ReadFontFromBytes returned error: %v", err)
	}

	if _, ok := font.chars[-5]; !ok {
		t.Fatalf("expected negative code tag to be loaded")
	}
}

func TestReadFontFromBytesRejectsInvalidPositiveCodeTags(t *testing.T) {
	_, err := ReadFontFromBytes([]byte(minimalFontWithCodeTag("0xD800", "T")))
	if err == nil {
		t.Fatal("expected invalid positive code tag to be rejected")
	}
}

func minimalFontWithCodeTag(codeTag string, glyph string) string {
	var b strings.Builder
	b.WriteString("flf2a$ 1 1 1 0 0\n")
	for i := 0; i < 102; i++ {
		b.WriteString("_@\n")
	}
	b.WriteString(codeTag)
	b.WriteByte('\n')
	b.WriteString(glyph)
	b.WriteString("@\n")

	return b.String()
}
