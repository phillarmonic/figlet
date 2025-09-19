# Figlet

This is a port of Figlet from C to the Go programming language.

This is a fork from upstream: https://github.com/zekroTJA/figlet

Figlet is a program that makes large letters out of ordinary text.

```
 ,          _   _         __        __         _     _ _
/|   |     | | | |        \ \      / /__  _ __| | __| | |
 |___|  _  | | | |  __     \ \ /\ / / _ \| '__| |/ _` | |
 |   |\|/  |/  |/  /  \_    \ V  V / (_) | |  | | (_| |_|
 |   |/|__/|__/|__/\__/      \_/\_/ \___/|_|  |_|\__,_(_)

```

For information about the original FIGlet, see [figlet.org](http://www.figlet.org/).

### Usage

```
figlet [ -lcrhR ] [ -f fontfile ] [ --gay ]
       [ --gradient-start color --gradient-end color ]
       [ -w outputwidth ] [ -m smushmode ]
       [ message ]
```

###### Options

`-h`
Shows help info: really just the usage info above plus the address of this page.

`-l, -c, -r`
These control the alignment of the output: left, center and right accordingly.

`-R`
Reverses the direction of text. So if the font specifies left-to-right, this will make it right-to-left, and vice versa.

`-f fontfile`
Specify a font to use. The fonts come from the "fonts" directory, in the same directory as the `figlet` program. You can see the available fonts with `figlet -list`.

`-w outputwidth`
FIGlet assumes an 80 character wide terminal. Use this to specify a different output width.

`-m smushmode`
Use a different "smush mode". Smush modes control how Figlet "smushes" together the big letters for output. This option is only really useful if you're making a font and need to experiment with the various settings—usually the font author has already specified the smush mode that works best with that font. You can find more information on smush modes in [figfont.txt](https://raw.github.com/lukesampson/figlet/master/figfont.txt), although this version of figfont.txt is written for the C version.

`-list`
Lists the available fonts, with a preview of each.

`--gay`
Rainbow colors (pride flag mode). Creates vibrant rainbow effects across the text.

`--gradient-start color` and `--gradient-end color`
Create a gradient from the start color to the end color across the text. Colors can be specified as:
- Hex colors: `#FF0000`, `#ff0000`, `FF0000`, `ff0000`
- RGB values: `rgb(255,0,0)`, `RGB(255,0,0)`
- Named colors: `red`, `blue`, `green`, `yellow`, `magenta`, `cyan`, `white`, `black`, `orange`, `purple`, `pink`, `lime`, `navy`, `teal`, `silver`, `gray`, `maroon`, `olive`

Both `--gradient-start` and `--gradient-end` must be specified together.

`message`
The message you want to print out. If you don't specify one, Figlet will go into interactive mode where it waits for you to enter a line of text and then prints it out in large letters. You can do this as many times as you like, and use Ctrl-C to quit.

### Color Examples

```bash
# Rainbow colors (pride flag mode)
figlet --gay "RAINBOW TEXT"

# Gradient from red to blue
figlet --gradient-start red --gradient-end blue "GRADIENT"

# Gradient with hex colors
figlet --gradient-start "#FF6B35" --gradient-end "#004E89" "HEX COLORS"

# Gradient with RGB values
figlet --gradient-start "rgb(255,0,255)" --gradient-end "rgb(0,255,0)" "RGB"

# Rainbow font listing
figlet --gay -list

# Gradient font listing
figlet --gradient-start orange --gradient-end purple -list
```

### Why did you port it?

I couldn't get [the C version](https://github.com/cmatsuoka/figlet) to build and run properly on Windows using MSYS. Rather than mess around with lots of things I don't understand, I decided this would be a good opportunity to learn Go instead.

Also, the original version of this program is over 20 years old, and the code shows it. The main loop has a comment that says:

    The following code is complex and thoroughly tested.
    Be careful when modifying!

I like to think that the Go version is a lot clearer, especially with a lot of the legacy options stripped out. Although I admit the Go code is not the best—this is my first time programming in Go. I'd appreciate pull requests that make it better.

### Differences from the original version

###### Control files

Control files aren't supported in this version. They seem like a legacy workaround for something that's not so much a problem any more. I've tested passing unicode characters directly to this version and it seems to work ok, when the font supports the character. Even if I haven't gotten it right, Go has excellent UTF8 support so it shouldn't be too hard to fix this in a way that doesn't involve the complexity of control files.

###### Newline handling

The original version has options for handling newlines, and I think it renders newlines as it receives them from input. This version just treats newlines as whitespace and won't print a new line by default. I might be wrong, but I think this is pretty much what you want in most cases anyway.

###### Unsupported options

These command-line options aren't supported in this version:

`-knopstvxDELNSWX`
Too complicated!

`-f fontdirectory`
This version tries to find the "fonts" directory in the same directory as the `figlet` executable. If you keep your fonts elsewhere, you can supply the `-p` flag.

`-C controlfile`
Control files aren't supported, for reasons given above.

`-I infocode`
Not supported

`-R`
This is supported, but it behaves differently in this version. In the original it meant "Right-to-left" print direction. In this version it means "Reverse" the print direction, as specified in the font file. Most times the font file is what you want, so this is mainly for testing and as a gimmick to confuse people.

### Use as a library

You can generate your own text at runtime using `figletlib`. The library supports **embedded fonts** (no external files needed) and **color output**.

#### **Basic Usage with Embedded Fonts**

```go
package main

import (
	"fmt"
	"os"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	// Use embedded fonts - no external files needed!
	loader := figletlib.NewEmbededLoader()
	
	font, err := loader.GetFontByName("standard")
	if err != nil {
		panic(err)
	}
	
	// Print to stdout
	figletlib.PrintMsg("Hello World!", font, 80, font.Settings(), "left")
}
```

#### **HTTP Server Example**

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/phillarmonic/figlet/figletlib"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Use embedded fonts - works anywhere!
	loader := figletlib.NewEmbededLoader()
	
	font, err := loader.GetFontByName("big")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Could not load font!")
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	figletlib.FPrintMsg(w, "Hello, Web!", font, 80, font.Settings(), "center")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

#### **Advanced Color Examples**

##### **Rainbow Colors (Gay Mode)**

```go
package main

import (
	"fmt"
	"os"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	font, err := loader.GetFontByName("big")
	if err != nil {
		panic(err)
	}
	
	// Rainbow colors (pride flag mode)
	rainbowConfig := figletlib.ColorConfig{
		Mode: figletlib.ColorModeRainbow,
	}
	
	fmt.Println("🏳️‍🌈 Rainbow Text Examples:")
	figletlib.PrintColoredMsg("PRIDE!", font, 80, font.Settings(), "center", rainbowConfig)
	
	// Different fonts with rainbow
	fonts := []string{"standard", "slant", "shadow", "bubble"}
	for _, fontName := range fonts {
		if f, err := loader.GetFontByName(fontName); err == nil {
			fmt.Printf("\n--- %s font ---\n", fontName)
			figletlib.PrintColoredMsg("RAINBOW", f, 80, f.Settings(), "left", rainbowConfig)
		}
	}
}
```

##### **Hex Color Gradients**

```go
package main

import (
	"fmt"
	"os"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	font, _ := loader.GetFontByName("big")
	
	// Hex colors with # prefix
	startColor, _ := figletlib.ParseColor("#FF6B35")  // Orange
	endColor, _ := figletlib.ParseColor("#004E89")    // Blue
	
	gradientConfig := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: startColor,
		EndColor:   endColor,
	}
	
	fmt.Println("🎨 Hex Color Gradient:")
	figletlib.PrintColoredMsg("HEX COLORS", font, 80, font.Settings(), "center", gradientConfig)
	
	// Hex colors without # prefix
	startColor2, _ := figletlib.ParseColor("FF1493")  // Deep Pink
	endColor2, _ := figletlib.ParseColor("00CED1")    // Dark Turquoise
	
	gradientConfig2 := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: startColor2,
		EndColor:   endColor2,
	}
	
	fmt.Println("\n🌺 Pink to Turquoise:")
	figletlib.PrintColoredMsg("BEAUTIFUL", font, 80, font.Settings(), "center", gradientConfig2)
}
```

##### **RGB Color Values**

```go
package main

import (
	"fmt"
	"os"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	font, _ := loader.GetFontByName("slant")
	
	// RGB format colors
	startColor, _ := figletlib.ParseColor("rgb(255, 0, 255)")  // Magenta
	endColor, _ := figletlib.ParseColor("rgb(0, 255, 0)")      // Lime
	
	gradientConfig := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: startColor,
		EndColor:   endColor,
	}
	
	fmt.Println("🔬 RGB Color Gradient:")
	figletlib.PrintColoredMsg("RGB VALUES", font, 80, font.Settings(), "center", gradientConfig)
	
	// Case insensitive RGB
	startColor2, _ := figletlib.ParseColor("RGB(255,165,0)")   // Orange
	endColor2, _ := figletlib.ParseColor("rgb(75,0,130)")      // Indigo
	
	gradientConfig2 := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: startColor2,
		EndColor:   endColor2,
	}
	
	fmt.Println("\n🌅 Orange to Indigo:")
	figletlib.PrintColoredMsg("SUNSET", font, 80, font.Settings(), "center", gradientConfig2)
}
```

##### **Direct RGB Struct Usage**

```go
package main

import (
	"fmt"
	"os"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	font, _ := loader.GetFontByName("shadow")
	
	// Create RGB colors directly
	neonPink := figletlib.RGB{R: 255, G: 20, B: 147}    // Hot Pink
	electricBlue := figletlib.RGB{R: 0, G: 191, B: 255}  // Deep Sky Blue
	
	gradientConfig := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: neonPink,
		EndColor:   electricBlue,
	}
	
	fmt.Println("⚡ Neon Colors:")
	figletlib.PrintColoredMsg("ELECTRIC", font, 80, font.Settings(), "center", gradientConfig)
	
	// Cyberpunk colors
	cyberGreen := figletlib.RGB{R: 0, G: 255, B: 65}
	cyberPurple := figletlib.RGB{R: 138, G: 43, B: 226}
	
	cyberConfig := figletlib.ColorConfig{
		Mode:       figletlib.ColorModeGradient,
		StartColor: cyberGreen,
		EndColor:   cyberPurple,
	}
	
	fmt.Println("\n🤖 Cyberpunk Style:")
	figletlib.PrintColoredMsg("CYBER", font, 80, font.Settings(), "center", cyberConfig)
}
```

##### **Color Interpolation & Custom Effects**

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	font, _ := loader.GetFontByName("big")
	
	// Create a custom color effect by interpolating between multiple colors
	colors := []figletlib.RGB{
		{255, 0, 0},    // Red
		{255, 165, 0},  // Orange  
		{255, 255, 0},  // Yellow
		{0, 255, 0},    // Green
		{0, 0, 255},    // Blue
		{75, 0, 130},   // Indigo
		{238, 130, 238}, // Violet
	}
	
	fmt.Println("🌈 Multi-Color Interpolation:")
	
	// Create multiple gradients to simulate rainbow effect
	for i := 0; i < len(colors)-1; i++ {
		gradientConfig := figletlib.ColorConfig{
			Mode:       figletlib.ColorModeGradient,
			StartColor: colors[i],
			EndColor:   colors[i+1],
		}
		
		word := fmt.Sprintf("COLOR%d", i+1)
		figletlib.PrintColoredMsg(word, font, 80, font.Settings(), "left", gradientConfig)
	}
}

// Custom color interpolation function
func interpolateColors(start, end figletlib.RGB, steps int) []figletlib.RGB {
	result := make([]figletlib.RGB, steps)
	for i := 0; i < steps; i++ {
		factor := float64(i) / float64(steps-1)
		result[i] = start.Interpolate(end, factor)
	}
	return result
}
```

##### **Colorful Web Server**

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	loader := figletlib.NewEmbededLoader()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		font, _ := loader.GetFontByName("big")
		
		// Rainbow welcome message
		rainbowConfig := figletlib.ColorConfig{
			Mode: figletlib.ColorModeRainbow,
		}
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "🌈 Welcome to the Colorful Figlet Server! 🌈")
		fmt.Fprintln(w)
		
		figletlib.FPrintColoredMsg(w, "WELCOME", font, 80, font.Settings(), "center", rainbowConfig)
		fmt.Fprintln(w)
	})
	
	http.HandleFunc("/gradient", func(w http.ResponseWriter, r *http.Request) {
		font, _ := loader.GetFontByName("slant")
		
		// Time-based gradient colors
		now := time.Now()
		hour := now.Hour()
		
		var startColor, endColor figletlib.RGB
		var message string
		
		switch {
		case hour >= 6 && hour < 12: // Morning
			startColor, _ = figletlib.ParseColor("#FFD700") // Gold
			endColor, _ = figletlib.ParseColor("#FF6347")   // Tomato
			message = "GOOD MORNING"
		case hour >= 12 && hour < 18: // Afternoon  
			startColor, _ = figletlib.ParseColor("#00BFFF") // Deep Sky Blue
			endColor, _ = figletlib.ParseColor("#32CD32")   // Lime Green
			message = "GOOD AFTERNOON"
		default: // Evening/Night
			startColor, _ = figletlib.ParseColor("#4B0082") // Indigo
			endColor, _ = figletlib.ParseColor("#8A2BE2")   // Blue Violet
			message = "GOOD EVENING"
		}
		
		gradientConfig := figletlib.ColorConfig{
			Mode:       figletlib.ColorModeGradient,
			StartColor: startColor,
			EndColor:   endColor,
		}
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "🕐 Time-based gradient (Hour: %d)\n\n", hour)
		figletlib.FPrintColoredMsg(w, message, font, 80, font.Settings(), "center", gradientConfig)
	})
	
	fmt.Println("🚀 Colorful server running on http://localhost:8080")
	fmt.Println("📍 Try http://localhost:8080/gradient for time-based colors!")
	http.ListenAndServe(":8080", nil)
}
```

#### **Combined Loader (External + Embedded)**

```go
package main

import (
	"github.com/phillarmonic/figlet/figletlib"
)

func main() {
	// Try external fonts first, fall back to embedded
	loader := figletlib.NewCombinedLoaderWithDir("./custom-fonts")
	
	// This will use custom font if available, or embedded font as fallback
	font, err := loader.GetFontByName("standard")
	if err != nil {
		panic(err)
	}
	
	figletlib.PrintMsg("Best of both worlds!", font, 80, font.Settings(), "left")
}
```

#### **Available Functions**

##### **Font Loading**
- `NewEmbededLoader()` - Use embedded fonts (recommended)
- `NewDirLoader(dir)` - Use fonts from directory  
- `NewCombinedLoaderWithDir(dir)` - Use external + embedded fonts

##### **Basic Text Output**
- `PrintMsg(msg, font, width, settings, align)` - Print to stdout
- `FPrintMsg(writer, msg, font, width, settings, align)` - Print to any io.Writer

##### **Colored Text Output**
- `PrintColoredMsg(msg, font, width, settings, align, colorConfig)` - Print with colors to stdout
- `FPrintColoredMsg(writer, msg, font, width, settings, align, colorConfig)` - Print with colors to any io.Writer

##### **Color Functions**
- `ParseColor(colorStr)` - Parse color strings (hex, rgb(), named colors)
- `RGB{R, G, B uint8}` - Create RGB color directly
- `(rgb RGB) ToANSI()` - Convert RGB to ANSI escape sequence
- `(start RGB) Interpolate(end RGB, factor float64)` - Interpolate between two colors
- `GetRainbowColor(position float64)` - Get rainbow color at position (0.0-1.0)
- `HSVtoRGB(h, s, v float64)` - Convert HSV to RGB

##### **Color Configuration**
- `ColorConfig{Mode, StartColor, EndColor}` - Color configuration struct
- `ColorModeNone` - No coloring
- `ColorModeGradient` - Gradient between two colors
- `ColorModeRainbow` - Rainbow/pride flag colors

##### **Supported Color Formats**
- **Hex**: `#FF0000`, `#ff0000`, `FF0000`, `ff0000`
- **RGB**: `rgb(255,0,0)`, `RGB(255,0,0)`
- **Named**: `red`, `blue`, `green`, `yellow`, `magenta`, `cyan`, `white`, `black`, `orange`, `purple`, `pink`, `lime`, `navy`, `teal`, `silver`, `gray`, `maroon`, `olive`
