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
	pkgName = "github.com/zekrotja/figlet"
)

var (
	//go:embed fonts
	embeddedFonts    embed.FS
	EmbeddedFonts, _ = fs.Sub(embeddedFonts, "fonts")
)

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
