package cli

import (
	"flag"

	"github.com/jcocozza/deck/internal/draw"
	"github.com/jcocozza/deck/internal/format"
	"github.com/jcocozza/deck/internal/render"
	"github.com/jcocozza/deck/internal/utils"
)

type deckcfg struct {
	Colorize bool
	NoScale  bool
	Theme    draw.Theme
}

func flags() ([]string, deckcfg) {
	theme := flag.String("theme", "", "which theme to use. must be in cfg unless using default options: empty = auto, 'default' = default color scheme")
	noScale := flag.Bool("no-scale", false, "don't auto scale - respect font sizes set in config")

	// TODO: figure this out
	// _ = flag.String("cfg", "", "config path - will check ~/.drawrc and ~/.config/drawrc by default")

	flag.Parse()

	var t draw.Theme
	switch *theme {
	case "":
		t = draw.DefaultTheme
	case "default":
		t = draw.DefaultColorTheme
	default:
		panic("unsupport theme type")
	}

	return flag.Args(), deckcfg{
		NoScale:  *noScale,
		Theme:    t,
	}
}

func Cli() error {
	args, cfg := flags()
	lines, err := utils.ReadFromStdinOrFiles(args)
	if err != nil {
		return err
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
		d = draw.NewDrawer(draw.Pretty, cfg.Theme)
	} else {
		d = draw.NewDrawer(draw.Auto, cfg.Theme)
	}
	return render.Render(slides, d)
}
