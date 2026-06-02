package figletlib

import (
	"bytes"
	"testing"
)

func TestFPrintLinesWithOptionsAppliesLeftMargin(t *testing.T) {
	var buf bytes.Buffer

	FPrintLinesWithOptions(&buf, []FigText{singleLineFigText("A")}, '$', 6, PrintOptions{
		Align:      "left",
		LeftMargin: 2,
	})

	if got := buf.String(); got != "  A\n" {
		t.Fatalf("unexpected output %q", got)
	}
}

func TestFPrintLinesWithOptionsCentersWithinRemainingWidth(t *testing.T) {
	var buf bytes.Buffer

	FPrintLinesWithOptions(&buf, []FigText{singleLineFigText("AB")}, '$', 10, PrintOptions{
		Align:      "center",
		LeftMargin: 2,
	})

	if got := buf.String(); got != "     AB\n" {
		t.Fatalf("unexpected output %q", got)
	}
}

func TestFPrintLinesWithOptionsIgnoresNegativeMargins(t *testing.T) {
	var buf bytes.Buffer

	FPrintLinesWithOptions(&buf, []FigText{singleLineFigText("A")}, '$', 4, PrintOptions{
		LeftMargin: -3,
	})

	if got := buf.String(); got != "A\n" {
		t.Fatalf("unexpected output %q", got)
	}
}

func singleLineFigText(line string) FigText {
	text := newFigText(1)
	text.art[0] = []rune(line)

	return *text
}
