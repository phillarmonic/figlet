package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/phillarmonic/figlet/figletlib"
)

const (
	defaultFont = "standard"
)

func printUsage() {
	fmt.Println("Usage: figlet [ -lcrhR ] [ -f fontfile ] [ -I infocode ]")
	fmt.Println("              [ -w outputwidth ] [ -m smushmode ]")
	fmt.Println("              [ message ]")
}

func printHelp() {
	printUsage()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("For more info see https://github.com/phillarmonic/figlet")
}

func printVersion(fontsSource string) {
	fmt.Println("Figlet version: go-1.0")
	fmt.Printf("Fonts: %v\n", fontsSource)
}

func printInfoCode(infocode int, infodata []string) {
	fmt.Println(infodata[infocode])
}

func listFonts(loader figletlib.FontLoader, fontsSource string, colorConfig figletlib.ColorConfig) {
	fmt.Printf("Fonts in %v:\n", fontsSource)
	fonts, _ := loader.FontNamesInDir()
	for _, fontname := range fonts {
		fmt.Printf("%v:\n", fontname)
		f, err := loader.GetFontByName(fontname)
		if err != nil {
			fmt.Println(err)
		}

		s := f.Settings()
		if colorConfig.Mode != figletlib.ColorModeNone {
			figletlib.PrintColoredMsg(fontname, f, 80, s, "left", colorConfig)
		} else {
			figletlib.PrintMsg(fontname, f, 80, s, "left")
		}
		fmt.Println()
	}
}

func main() {
	// options
	fontname := flag.String("f", defaultFont, "name of font to use")
	reverse := flag.Bool("R", false, "reverse output")
	alignRight := flag.Bool("r", false, "right-align output")
	alignCenter := flag.Bool("c", false, "center-align output")
	outputWidth := flag.Int("w", 80, "output width")
	list := flag.Bool("list", false, "list available fonts")
	help := flag.Bool("h", false, "show help")
	version := flag.Bool("v", false, "show version info")
	fontsDirectory := flag.String("d", "", "fonts directory")
	infoCode := flag.Int("I", -1, "infocode")
	infoCode2 := flag.Bool("I2", false, "show default font directory")
	infoCode3 := flag.Bool("I3", false, "show default font")
	infoCode4 := flag.Bool("I4", false, "show output width")
	infoCode5 := flag.Bool("I5", false, "show supported font formats")

	// color options
	gradientStart := flag.String("gradient-start", "", "start color for gradient (hex, rgb(), or named color)")
	gradientEnd := flag.String("gradient-end", "", "end color for gradient (hex, rgb(), or named color)")
	rainbow := flag.Bool("gay", false, "rainbow colors (pride flag mode)")
	flag.Parse()
	var loader figletlib.FontLoader
	fontsdir := *fontsDirectory
	var fontsSource string

	if fontsdir == "" {
		// Try to auto-detect fonts directory
		detectedDir := figletlib.GuessFontsDirectory()
		if detectedDir != "" {
			// Use combined loader: detected directory + embedded fonts
			loader = figletlib.NewCombinedLoaderWithDir(detectedDir)
			fontsSource = fmt.Sprintf("%s + embedded", detectedDir)
		} else {
			// Use only embedded fonts
			loader = figletlib.NewEmbededLoader()
			fontsSource = "embedded"
		}
	} else {
		// Use combined loader: specified directory + embedded fonts
		loader = figletlib.NewCombinedLoaderWithDir(fontsdir)
		fontsSource = fmt.Sprintf("%s + embedded", fontsdir)
	}

	// Configure colors early so we can use them in list mode
	var colorConfig figletlib.ColorConfig
	colorConfig.Mode = figletlib.ColorModeNone

	if *rainbow {
		colorConfig.Mode = figletlib.ColorModeRainbow
	} else if *gradientStart != "" && *gradientEnd != "" {
		startColor, err := figletlib.ParseColor(*gradientStart)
		if err != nil {
			fmt.Printf("ERROR: invalid start color '%s': %v\n", *gradientStart, err)
			os.Exit(1)
		}

		endColor, err := figletlib.ParseColor(*gradientEnd)
		if err != nil {
			fmt.Printf("ERROR: invalid end color '%s': %v\n", *gradientEnd, err)
			os.Exit(1)
		}

		colorConfig.Mode = figletlib.ColorModeGradient
		colorConfig.StartColor = startColor
		colorConfig.EndColor = endColor
	} else if *gradientStart != "" || *gradientEnd != "" {
		fmt.Println("ERROR: both --gradient-start and --gradient-end must be specified for gradient mode")
		os.Exit(1)
	}

	if *list {
		listFonts(loader, fontsSource, colorConfig)
		os.Exit(0)
	}

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *version {
		printVersion(fontsSource)
		os.Exit(0)
	}

	var align string
	if *alignRight {
		align = "right"
	} else if *alignCenter {
		align = "center"
	}

	f, err := loader.GetFontByName(*fontname)
	if err != nil {
		fmt.Println("ERROR: couldn't find font", *fontname, "in", fontsSource)
		os.Exit(1)
	}

	msg := strings.Join(flag.Args(), " ")

	s := f.Settings()
	if *reverse {
		s.SetRtoL(true)
	}

	ic := *infoCode

	if *infoCode2 {
		ic = 2
	} else if *infoCode3 {
		ic = 3
	} else if *infoCode4 {
		ic = 4
	} else if *infoCode5 {
		ic = 5
	}

	if ic > 1 && ic < 6 {
		outputWidthString := strconv.Itoa(*outputWidth)
		infoData := []string{2: fontsSource, 3: *fontname, 4: outputWidthString, 5: "flf2"}
		printInfoCode(ic, infoData)
		os.Exit(0)
	} else if ic != -1 {
		fmt.Println("ERROR: invalid infocode", ic)
		os.Exit(1)
	}

	maxwidth := *outputWidth
	if msg == "" {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err = reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					os.Exit(0)
				}
				msg = ""
			}
			if colorConfig.Mode != figletlib.ColorModeNone {
				figletlib.PrintColoredMsg(msg, f, maxwidth, s, align, colorConfig)
			} else {
				figletlib.PrintMsg(msg, f, maxwidth, s, align)
			}
		}
	}
	if colorConfig.Mode != figletlib.ColorModeNone {
		figletlib.PrintColoredMsg(msg, f, maxwidth, s, align, colorConfig)
	} else {
		figletlib.PrintMsg(msg, f, maxwidth, s, align)
	}
}
