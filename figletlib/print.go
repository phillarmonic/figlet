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

type PrintOptions struct {
	Align      string
	LeftMargin int
}

func normalizePrintOptions(options PrintOptions) PrintOptions {
	if options.LeftMargin < 0 {
		options.LeftMargin = 0
	}

	return options
}

func contentWidth(maxwidth int, options PrintOptions) int {
	width := maxwidth - options.LeftMargin
	if width < 0 {
		return 0
	}

	return width
}

func leftPadding(linelen int, maxwidth int, options PrintOptions) int {
	padding := options.LeftMargin
	width := contentWidth(maxwidth, options)
	if width <= linelen {
		return padding
	}

	switch options.Align {
	case alignRight:
		padding += width - linelen
	case "center":
		padding += (width - linelen) / 2
	}

	return padding
}

func writeLeftPadding(w io.Writer, linelen int, maxwidth int, options PrintOptions) {
	if padding := leftPadding(linelen, maxwidth, options); padding > 0 {
		_, _ = fmt.Fprint(w, strings.Repeat(" ", padding))
	}
}

func FPrintLines(w io.Writer, lines []FigText, hardblank rune, maxwidth int, align string) {
	FPrintLinesWithOptions(w, lines, hardblank, maxwidth, PrintOptions{Align: align})
}

func FPrintLinesWithOptions(w io.Writer, lines []FigText, hardblank rune, maxwidth int, options PrintOptions) {
	options = normalizePrintOptions(options)
	width := contentWidth(maxwidth, options)
	for _, line := range lines {
		for _, subline := range line.Art() {
			writeLeftPadding(w, len(subline), maxwidth, options)
			for _, outchar := range subline {
				if outchar == hardblank {
					outchar = ' '
				}
				_, _ = fmt.Fprintf(w, "%c", outchar)
			}
			if len(subline) < width && options.Align != alignRight {
				_, _ = fmt.Fprintln(w)
			}
		}
	}
}

func PrintLines(lines []FigText, hardblank rune, maxwidth int, align string) {
	FPrintLines(os.Stdout, lines, hardblank, maxwidth, align)
}

func PrintLinesWithOptions(lines []FigText, hardblank rune, maxwidth int, options PrintOptions) {
	FPrintLinesWithOptions(os.Stdout, lines, hardblank, maxwidth, options)
}

func FPrintMsg(w io.Writer, msg string, f *Font, maxwidth int, s Settings, align string) {
	FPrintMsgWithOptions(w, msg, f, maxwidth, s, PrintOptions{Align: align})
}

func FPrintMsgWithOptions(w io.Writer, msg string, f *Font, maxwidth int, s Settings, options PrintOptions) {
	lines := GetLines(msg, f, maxwidth, s)
	FPrintLinesWithOptions(w, lines, s.HardBlank(), maxwidth, options)
}

func PrintMsg(msg string, f *Font, maxwidth int, s Settings, align string) {
	FPrintMsg(os.Stdout, msg, f, maxwidth, s, align)
}

func PrintMsgWithOptions(msg string, f *Font, maxwidth int, s Settings, options PrintOptions) {
	FPrintMsgWithOptions(os.Stdout, msg, f, maxwidth, s, options)
}

func SprintMsg(msg string, f *Font, maxwidth int, s Settings, align string) string {
	return SprintMsgWithOptions(msg, f, maxwidth, s, PrintOptions{Align: align})
}

func SprintMsgWithOptions(msg string, f *Font, maxwidth int, s Settings, options PrintOptions) string {
	buf := bytes.NewBufferString("")
	FPrintMsgWithOptions(buf, msg, f, maxwidth, s, options)

	return buf.String()
}

// FPrintColoredLines prints lines with color support.
func FPrintColoredLines(w io.Writer, lines []FigText, hardblank rune, maxwidth int, align string, colorConfig ColorConfig) {
	FPrintColoredLinesWithOptions(w, lines, hardblank, maxwidth, PrintOptions{Align: align}, colorConfig)
}

func FPrintColoredLinesWithOptions(w io.Writer, lines []FigText, hardblank rune, maxwidth int, options PrintOptions, colorConfig ColorConfig) {
	options = normalizePrintOptions(options)
	width := contentWidth(maxwidth, options)
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
			writeLeftPadding(w, len(subline), maxwidth, options)
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
			if len(subline) < width && options.Align != alignRight {
				_, _ = fmt.Fprintln(w)
			}
		}
	}
}

// PrintColoredLines prints lines with color support to stdout.
func PrintColoredLines(lines []FigText, hardblank rune, maxwidth int, align string, colorConfig ColorConfig) {
	FPrintColoredLines(os.Stdout, lines, hardblank, maxwidth, align, colorConfig)
}

func PrintColoredLinesWithOptions(lines []FigText, hardblank rune, maxwidth int, options PrintOptions, colorConfig ColorConfig) {
	FPrintColoredLinesWithOptions(os.Stdout, lines, hardblank, maxwidth, options, colorConfig)
}

// FPrintColoredMsg prints a message with color support.
func FPrintColoredMsg(w io.Writer, msg string, f *Font, maxwidth int, s Settings, align string, colorConfig ColorConfig) {
	FPrintColoredMsgWithOptions(w, msg, f, maxwidth, s, PrintOptions{Align: align}, colorConfig)
}

func FPrintColoredMsgWithOptions(w io.Writer, msg string, f *Font, maxwidth int, s Settings, options PrintOptions, colorConfig ColorConfig) {
	lines := GetLines(msg, f, maxwidth, s)
	FPrintColoredLinesWithOptions(w, lines, s.HardBlank(), maxwidth, options, colorConfig)
}

// PrintColoredMsg prints a message with color support to stdout.
func PrintColoredMsg(msg string, f *Font, maxwidth int, s Settings, align string, colorConfig ColorConfig) {
	FPrintColoredMsg(os.Stdout, msg, f, maxwidth, s, align, colorConfig)
}

func PrintColoredMsgWithOptions(msg string, f *Font, maxwidth int, s Settings, options PrintOptions, colorConfig ColorConfig) {
	FPrintColoredMsgWithOptions(os.Stdout, msg, f, maxwidth, s, options, colorConfig)
}
