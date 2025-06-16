package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/deck/internal/conf"
	"github.com/jcocozza/deck/internal/draw"
	"github.com/jcocozza/deck/internal/format"
	"github.com/jcocozza/deck/internal/render"
	"github.com/jcocozza/deck/internal/utils"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] [files...] \n", os.Args[0])
	flag.PrintDefaults()
}

type deckcfg struct {
	Colorize bool
	NoScale  bool
	Theme    string
}

func flags() ([]string, deckcfg) {
	theme := flag.String("theme", "", "which theme to use. must be in cfg unless using default options: empty = auto, 'default' = default color scheme")
	noScale := flag.Bool("no-scale", false, "don't auto scale - respect font sizes set in the current theme")
	colorize := flag.Bool("colorize", false, "fallback to default colors theme (fallback only happens when custom theme is missing parts)")

	flag.Parse()

	return flag.Args(), deckcfg{
		NoScale:  *noScale,
		Theme:    *theme,
		Colorize: *colorize,
	}
}

func Cli() {
	flag.Usage = usage
	nocfgFound := false
	c, err := conf.ReadConfig()
	if errors.Is(err, os.ErrNotExist) {
		nocfgFound = true
	} else if err != nil {
		ExitError(err)
	}
	args, cfgArgs := flags()
	lines, err := utils.ReadFromStdinOrFiles(args)
	if err != nil {
		ExitError(err)
	}

	var theme draw.Theme
	switch cfgArgs.Theme {
	case "":
		theme = draw.DefaultTheme
	case "default":
		theme = draw.DefaultColorTheme
	default:
		ctheme, ok := c.Themes[cfgArgs.Theme]
		if !ok || nocfgFound {
			Warn(fmt.Sprintf("theme %s not found. using default", cfgArgs.Theme))
			theme = draw.DefaultTheme
		} else {
			if cfgArgs.Colorize {
				theme, err = conf.LinkTheme(ctheme, draw.DefaultColorTheme)
				if err != nil {
					Warn(fmt.Sprintf("error linking theme %s. %s", cfgArgs.Theme, err.Error()))
				}
			} else {
				theme, err = conf.LinkTheme(ctheme, draw.DefaultTheme)
				if err != nil {
					Warn(fmt.Sprintf("error linking theme %s. %s", cfgArgs.Theme, err.Error()))
				}
			}
		}
	}

	lexer := format.NewLexer()
	parser := format.NewParser()

	lexLines := lexer.Lex(lines)
	slides, err := parser.Parse(lexLines)
	if err != nil {
		ExitError(err)
	}

	var d draw.Drawer
	if cfgArgs.NoScale {
		d = draw.NewDrawer(draw.Pretty, theme)
	} else {
		d = draw.NewDrawer(draw.Auto, theme)
	}
	err = render.Render(slides, d)
	if err != nil {
		ExitError(err)
	}
}
