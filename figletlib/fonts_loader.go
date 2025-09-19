package figletlib

import (
	"embed"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	pkgName = "github.com/phillarmonic/figlet"
)

var (
	//go:embed fonts
	embeddedFonts    embed.FS
	EmbeddedFonts, _ = fs.Sub(embeddedFonts, "fonts")
)

// FontLoader interface for loading fonts
type FontLoader interface {
	FontNamesInDir() ([]string, error)
	GetFontByName(name string) (*Font, error)
}

type Loader struct {
	fsys fs.FS
}

func NewLoader(fsys fs.FS) *Loader {
	return &Loader{fsys: fsys}
}

func NewDirLoader(dir string) *Loader {
	return NewLoader(os.DirFS(dir))
}

func NewEmbededLoader() *Loader {
	return NewLoader(EmbeddedFonts)
}

// CombinedLoader can access fonts from multiple sources
type CombinedLoader struct {
	loaders []*Loader
}

func NewCombinedLoader(loaders ...*Loader) *CombinedLoader {
	return &CombinedLoader{loaders: loaders}
}

// NewCombinedLoaderWithDir creates a combined loader with external directory and embedded fonts
func NewCombinedLoaderWithDir(dir string) *CombinedLoader {
	loaders := []*Loader{}

	// Add external directory loader if directory exists and has fonts
	if dir != "" {
		if fontsGlob := filepath.Join(dir, "*.flf"); fontsGlob != "" {
			if matches, err := filepath.Glob(fontsGlob); err == nil && len(matches) > 0 {
				loaders = append(loaders, NewDirLoader(dir))
			}
		}
	}

	// Always add embedded fonts as fallback
	loaders = append(loaders, NewEmbededLoader())

	return NewCombinedLoader(loaders...)
}

func (cl *CombinedLoader) FontNamesInDir() ([]string, error) {
	fontNamesMap := make(map[string]bool)
	var allFontNames []string

	for _, loader := range cl.loaders {
		names, err := loader.FontNamesInDir()
		if err != nil {
			continue // Skip loaders that fail
		}

		for _, name := range names {
			if !fontNamesMap[name] {
				fontNamesMap[name] = true
				allFontNames = append(allFontNames, name)
			}
		}
	}

	return allFontNames, nil
}

func (cl *CombinedLoader) GetFontByName(name string) (*Font, error) {
	var lastErr error

	for _, loader := range cl.loaders {
		font, err := loader.GetFontByName(name)
		if err == nil {
			return font, nil
		}
		lastErr = err
	}

	return nil, lastErr
}

func GuessFontsDirectory() string {
	bin := os.Args[0]
	if !filepath.IsAbs(bin) {
		maybeBin, err := filepath.Abs(bin)
		if err == nil {
			bin = maybeBin
		}
	}

	// try <bindir>
	bindir := filepath.Dir(bin)
	dirsToTry := []string{
		filepath.Join(bindir, "figletlib", "fonts"),
		filepath.Join(bindir, "fonts"),
	}

	// try src directory
	ctx := build.Default
	if p, err := ctx.Import(pkgName, "", build.FindOnly); err == nil {
		dirsToTry = append(dirsToTry, filepath.Join(p.Dir, "figletlib", "fonts"))
		dirsToTry = append(dirsToTry, filepath.Join(p.Dir, "fonts"))
	}

	for _, fontsDir := range dirsToTry {
		fontsGlob := filepath.Join(fontsDir, "*.flf")
		matches, err := filepath.Glob(fontsGlob)
		if err == nil && len(matches) > 0 {
			return fontsDir
		}
	}

	return ""
}

func (t *Loader) FontNamesInDir() ([]string, error) {
	matches, err := fs.Glob(t.fsys, "*.flf")
	if err != nil {
		return nil, err
	}

	fontNames := make([]string, 0)
	for _, filename := range matches {
		base := filepath.Base(filename)
		fontNames = append(fontNames, strings.TrimSuffix(base, ".flf"))
	}

	return fontNames, nil
}

func (t *Loader) GetFontByName(name string) (*Font, error) {
	if !strings.HasSuffix(name, ".flf") {
		name += ".flf"
	}

	return t.ReadFont(name)
}
