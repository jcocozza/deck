package conf

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcocozza/deck/internal/draw"
	"github.com/jcocozza/deck/internal/utils"
)

type Config struct {
	Themes map[string]Theme `json:"themes"`
}

const deckrc = "deckrc.json"

func ReadConfig() (Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	path := fmt.Sprintf("%s/%s",dir, deckrc)

	cfgBytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func linkThemeElement(te *themeElement, defaultThemeElement draw.ThemeElement) (draw.ThemeElement, error) {
	dt := draw.ThemeElement{}
	if te == nil {
		return defaultThemeElement, nil
	}

	if te.Size == nil {
		dt.Size = defaultThemeElement.Size
	} else {
		dt.Size = *te.Size
	}
	if te.Color == nil || *te.Color == "" {
		dt.Color = defaultThemeElement.Color
	} else {
		c, err := utils.ParseHexColor(*te.Color)
		if err != nil {
			return defaultThemeElement, err
		}
		dt.Color = c
	}
	return dt, nil
}

func LinkTheme(t Theme, defaultTheme draw.Theme) (draw.Theme, error) {
	dt := draw.Theme{}
	if t.Background == nil || *t.Background == "" {
		dt.Background = defaultTheme.Background
	} else {
		c, err := utils.ParseHexColor(*t.Background)
		if err != nil {
			return defaultTheme, err
		}
		dt.Background = c
	}
	var err error
	dt.Header, err = linkThemeElement(t.Header, defaultTheme.Header)
	dt.SubHeader, err = linkThemeElement(t.SubHeader, defaultTheme.SubHeader)
	dt.SubSubHeader, err = linkThemeElement(t.SubSubHeader, defaultTheme.SubSubHeader)
	dt.Link, err = linkThemeElement(t.Link, defaultTheme.Link)
	dt.Default, err = linkThemeElement(t.Default, defaultTheme.Default)
	if err != nil {
		return defaultTheme, err
	}
	return dt, nil
}
