package figletlib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	alignRight = "right"
)

func FPrintLines(w io.Writer, lines []FigText, hardblank rune, maxwidth int, align string) {
	padleft := func(linelen int) {
		switch align {
		case alignRight:
			_, _ = fmt.Fprint(w, strings.Repeat(" ", maxwidth-linelen))
		case "center":
			_, _ = fmt.Fprint(w, strings.Repeat(" ", (maxwidth-linelen)/2))
		}
	}

	for _, line := range lines {
		for _, subline := range line.Art() {
			padleft(len(subline))
			for _, outchar := range subline {
				if outchar == hardblank {
					outchar = ' '
				}
				_, _ = fmt.Fprintf(w, "%c", outchar)
			}
			if len(subline) < maxwidth && align != alignRight {
				_, _ = fmt.Fprintln(w)
			}
		}
	}
}

func PrintLines(lines []FigText, hardblank rune, maxwidth int, align string) {
	FPrintLines(os.Stdout, lines, hardblank, maxwidth, align)
}

func FPrintMsg(w io.Writer, msg string, f *Font, maxwidth int, s Settings, align string) {
	lines := GetLines(msg, f, maxwidth, s)
	FPrintLines(w, lines, s.HardBlank(), maxwidth, align)
}

func PrintMsg(msg string, f *Font, maxwidth int, s Settings, align string) {
	FPrintMsg(os.Stdout, msg, f, maxwidth, s, align)
}

func SprintMsg(msg string, f *Font, maxwidth int, s Settings, align string) string {
	buf := bytes.NewBufferString("")
	FPrintMsg(buf, msg, f, maxwidth, s, align)

	return buf.String()
}

// FPrintColoredLines prints lines with color support.
func FPrintColoredLines(w io.Writer, lines []FigText, hardblank rune, maxwidth int, align string, colorConfig ColorConfig) {
	padleft := func(linelen int) {
		switch align {
		case alignRight:
			_, _ = fmt.Fprint(w, strings.Repeat(" ", maxwidth-linelen))
		case "center":
			_, _ = fmt.Fprint(w, strings.Repeat(" ", (maxwidth-linelen)/2))
		}
	}

	for _, line := range lines {
		art := line.Art()
		if len(art) == 0 {
			continue
		}

		// Calculate total width for gradient calculation
		totalWidth := 0
		if len(art) > 0 {
			totalWidth = len(art[0])
		}

		for _, subline := range art {
			padleft(len(subline))
			for i, outchar := range subline {
				if outchar == hardblank {
					outchar = ' '
				}

				// Calculate position for gradient (0.0 to 1.0)
				position := 0.0
				if totalWidth > 1 {
					position = float64(i) / float64(totalWidth-1)
				}

				// Apply color if character is not a space
				if outchar != ' ' && colorConfig.Mode != ColorModeNone {
					coloredChar := ApplyColor(outchar, position, totalWidth, colorConfig)
					_, _ = fmt.Fprint(w, coloredChar)
				} else {
					_, _ = fmt.Fprintf(w, "%c", outchar)
				}
			}
			if len(subline) < maxwidth && align != alignRight {
				_, _ = fmt.Fprintln(w)
			}
		}
	}
}

// PrintColoredLines prints lines with color support to stdout.
func PrintColoredLines(lines []FigText, hardblank rune, maxwidth int, align string, colorConfig ColorConfig) {
	FPrintColoredLines(os.Stdout, lines, hardblank, maxwidth, align, colorConfig)
}

// FPrintColoredMsg prints a message with color support.
func FPrintColoredMsg(w io.Writer, msg string, f *Font, maxwidth int, s Settings, align string, colorConfig ColorConfig) {
	lines := GetLines(msg, f, maxwidth, s)
	FPrintColoredLines(w, lines, s.HardBlank(), maxwidth, align, colorConfig)
}

// PrintColoredMsg prints a message with color support to stdout.
func PrintColoredMsg(msg string, f *Font, maxwidth int, s Settings, align string, colorConfig ColorConfig) {
	FPrintColoredMsg(os.Stdout, msg, f, maxwidth, s, align, colorConfig)
}
