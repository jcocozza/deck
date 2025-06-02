package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/deck/internal/conf"
	"github.com/jcocozza/deck/internal/draw"
	"github.com/jcocozza/deck/internal/format"
	"github.com/jcocozza/deck/internal/render"
	"github.com/jcocozza/deck/internal/utils"
)

type deckcfg struct {
	Colorize bool
	NoScale  bool
	Theme    string
}

func flags() ([]string, deckcfg) {
	theme := flag.String("theme", "", "which theme to use. must be in cfg unless using default options: empty = auto, 'default' = default color scheme")
	noScale := flag.Bool("no-scale", false, "don't auto scale - respect font sizes set in the current theme")
	colorize := flag.Bool("colorize", false, "this flag will set the default theme to include more colors")

	flag.Parse()

	return flag.Args(), deckcfg{
		NoScale:  *noScale,
		Theme:    *theme,
		Colorize: *colorize,
	}
}

func Cli() error {
	dir, err := os.UserHomeDir()
	if err != nil { return err }
	c, err := conf.ReadConfig(fmt.Sprintf("%s/deckrc.json",dir))
	if err != nil {
		return err
	}
	args, cfg := flags()
	lines, err := utils.ReadFromStdinOrFiles(args)
	if err != nil {
		return err
	}

	var theme draw.Theme
	switch cfg.Theme {
	case "":
		theme = draw.DefaultTheme
	case "default":
		theme = draw.DefaultColorTheme
	default:
		ctheme, ok := c.Themes[cfg.Theme]
		if !ok {
			theme = draw.DefaultTheme
		} else {
			if cfg.Colorize {
				theme, err = conf.LinkTheme(ctheme, draw.DefaultColorTheme)
				if err != nil { fmt.Println(err.Error())}
			} else {
				theme, err = conf.LinkTheme(ctheme, draw.DefaultTheme)
				if err != nil { fmt.Println(err.Error())}
			}
		}
	}

	lexer := format.NewLexer()
	parser := format.NewParser()

	lexLines := lexer.Lex(lines)
	// slides := slide.TextSlides()
	slides, err := parser.Parse(lexLines)
	if err != nil {
		return err
	}

	var d draw.Drawer
	if cfg.NoScale {
		d = draw.NewDrawer(draw.Pretty, theme)
	} else {
		d = draw.NewDrawer(draw.Auto, theme)
	}
	return render.Render(slides, d)
}
