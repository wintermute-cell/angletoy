package colorscheme

import (
	"gorl/fw/core/logging"
	"image/color"
	"strconv"
	"strings"
)

type Hex string

var Colorscheme = struct {
	Color01 Hex
	Color02 Hex
	Color03 Hex
	Color04 Hex
	Color05 Hex
	Color06 Hex
	Color07 Hex
	Color08 Hex
	Color09 Hex
	Color10 Hex
	Color11 Hex
	Color12 Hex
	Color13 Hex
	Color14 Hex
	Color15 Hex
	Color16 Hex
}{
	"#8c8fae", // 1
	"#584563", // 2
	"#3e2137", // 3
	"#9a6348", // 4
	"#d79b7d", // 5
	"#f5edba", // 6
	"#c0c741", // 7
	"#647d34", // 8
	"#e4943a", // 9
	"#9d303b", // 10
	"#d26471", // 11
	"#70377f", // 12
	"#7ec4c1", // 13
	"#34859d", // 14
	"#17434b", // 15
	"#1f0e1c", // 16
}

func (h Hex) ToRGBA() color.RGBA {
	rgba, err := hex2RGBA(h)
	if err != nil {
		logging.Error("Failed to convert hex to RGBA: %v", err)
		rgba = color.RGBA{ // A strong pink color for errors
			R: 255,
			G: 51,
			B: 194,
			A: 255,
		}
	}
	return rgba
}

func hex2RGBA(hex Hex) (color.RGBA, error) {
	var rgb color.RGBA
	values, err := strconv.ParseUint(strings.Trim(string(hex), "#"), 16, 32)

	if err != nil {
		return color.RGBA{}, err
	}

	rgb = color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
		A: 0xFF,
	}

	return rgb, nil
}
