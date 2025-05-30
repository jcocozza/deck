package cli

import (
	"flag"

	"github.com/jcocozza/deck/internal/draw"
	"github.com/jcocozza/deck/internal/format"
	"github.com/jcocozza/deck/internal/render"
	"github.com/jcocozza/deck/internal/utils"
)

type deckcfg struct {
	Pretty bool
	Theme draw.Theme
}

func flags() ([]string, deckcfg) {
	pretty := flag.Bool("pretty", false, "use pretty output")
	// TODO: figure this out
	_ = flag.String("cfg", "", "config path - will check ~/.drawrc and ~/.config/drawrc by default")

	flag.Parse()

	return flag.Args(), deckcfg{
		Pretty: *pretty,
		Theme: draw.DefaultTheme,
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
	if cfg.Pretty {
		d = draw.NewDrawer(draw.Pretty, cfg.Theme)
	} else {
		d = draw.NewDrawer(draw.Auto, cfg.Theme)
	}
	return render.Render(slides, d)
}
