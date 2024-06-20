package gui

import (
	"gorl/fw/core/logging"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strconv"
	"strings"
)

/*
parseStyleDefinition converts a string of style definitions into a map representation.

Input format: "property:value|property2:value2|...".

Note:

- If a value itself contains colons, the entire value after the first colon is
considered as the value. For example, "url:http:\/\/example.com" will be parsed
as {"url": "http:\/\/example.com"}.

- Extra spaces between properties, colons, and values are trimmed in the output.
*/
func parseStyleDef(style_definition string) map[string]any {
	if style_definition == "" {
		return make(map[string]any)
	}

	// remove potential leading/trailing pipe symbol
	style_definition = strings.Trim(style_definition, "|")

	pairs := strings.Split(style_definition, "|")
	style_map := make(map[string]any)

	// Define conversion functions
	converters := make(map[string]func(string) any)

	converters["color"] = func(value string) any {
		rgba := strings.Split(value, ",")
		if len(rgba) != 4 {
			logging.Warning("Invalid color format for value: %v", value)
			return nil
		}

		var result [4]uint8
		for i, v := range rgba {
			f, err := strconv.Atoi(v) // directly parse to int
			if err != nil || f < 0 || f > 255 {
				logging.Warning("Invalid color component in value: %v", value)
				return nil
			}
			result[i] = uint8(f)
		}
		return rl.NewColor(result[0], result[1], result[2], result[3])
	}

	converters["background"] = converters["color"]

	converters["background-hover"] = converters["color"]

	converters["background-pressed"] = converters["color"]

	converters["font"] = func(value string) any {
		return value // already a string
	}

	converters["font-scale"] = func(value string) any {
		s, err := strconv.ParseFloat(value, 32)
		if err != nil {
			logging.Warning("Invalid float value: %v", value)
		}
		return float32(s)
	}

	converters["debug"] = func(value string) any {
		b, err := strconv.ParseBool(value)
		if err != nil {
			logging.Warning("Invalid bool value: %v", value)
		}
		return b
	}
	// add more converters here as needed

	// iterate over all property:value pairs, and try to apply the appropriate
	// converter function.
	for _, pair := range pairs {
		if pair == "" {
			logging.Warning("Found empty pair in gui styledef!")
			continue
		}
		pv := strings.SplitN(pair, ":", 2)
		if pv[0] == "" || len(pv) != 2 || pv[1] == "" {
			logging.Warning("Found empty property or value in gui styledef: %v", pair)
			continue
		}

		// Trim spaces in case there are any around the property or value
		prop := strings.TrimSpace(pv[0])
		val := strings.TrimSpace(pv[1])

		// Convert to appropriate data type
		if converter, ok := converters[prop]; ok {
			style_map[prop] = converter(val)
		} else {
			logging.Warning("Unknown property in gui styledef: %v", prop)
		}
	}

	return style_map
}
