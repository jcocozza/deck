package utils

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

func ParseHexColor(s string) (color.Color, error) {
	// Strip leading "#", if present
	s = strings.TrimPrefix(s, "#")

	var r, g, b uint8
	var err error

	switch len(s) {
	case 6:
		// RRGGBB
		rVal, err1 := strconv.ParseUint(s[0:2], 16, 8)
		gVal, err2 := strconv.ParseUint(s[2:4], 16, 8)
		bVal, err3 := strconv.ParseUint(s[4:6], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			err = fmt.Errorf("invalid hex color %v, %v, %v", err1, err2, err3)
		}
		r, g, b = uint8(rVal), uint8(gVal), uint8(bVal)
	default:
		err = fmt.Errorf("unsupported hex format: %s", s)
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, err
}
