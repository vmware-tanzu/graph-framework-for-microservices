package nexus_compare

// This part is a fork of lib https://github.com/homeport/dyff

import (
	"fmt"
	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/ytbx"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	additionGreen      = color("#58BF38")
	modificationYellow = color("#C7C43F")
	removalRed         = color("#B9311B")
)

func yamlStringInGreenishColors(input interface{}) (string, error) {
	return neat.NewOutputProcessor(true, true, &map[string]colorful.Color{
		"keyColor":           bunt.Green,
		"indentLineColor":    {R: 0, G: 0.2, B: 0},
		"scalarDefaultColor": bunt.LimeGreen,
		"boolColor":          bunt.LimeGreen,
		"floatColor":         bunt.LimeGreen,
		"intColor":           bunt.LimeGreen,
		"multiLineTextColor": bunt.OliveDrab,
		"nullColor":          bunt.Olive,
		"emptyStructures":    bunt.DarkOliveGreen,
		"dashColor":          bunt.Green,
	}).ToYAML(input)
}

func yamlStringInRedishColors(input interface{}) (string, error) {
	return neat.NewOutputProcessor(true, true, &map[string]colorful.Color{
		"keyColor":           bunt.FireBrick,
		"indentLineColor":    {R: 0.2, G: 0, B: 0},
		"scalarDefaultColor": bunt.LightCoral,
		"boolColor":          bunt.LightCoral,
		"floatColor":         bunt.LightCoral,
		"intColor":           bunt.LightCoral,
		"multiLineTextColor": bunt.DarkSalmon,
		"nullColor":          bunt.Salmon,
		"emptyStructures":    bunt.LightSalmon,
		"dashColor":          bunt.FireBrick,
	}).ToYAML(input)
}

func colored(color colorful.Color, format string, a ...interface{}) string {
	return bunt.Style(
		render(format, a...),
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}

func color(hex string) colorful.Color {
	color, _ := colorful.Hex(hex)
	return color
}

func render(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}

	return fmt.Sprintf(format, a...)
}

func green(format string, a ...interface{}) string {
	return colored(additionGreen, render(format, a...))
}

func red(format string, a ...interface{}) string {
	return colored(removalRed, render(format, a...))
}

func yellow(format string, a ...interface{}) string {
	return colored(modificationYellow, render(format, a...))
}

func lightgreen(format string, a ...interface{}) string {
	return colored(bunt.LightGreen, render(format, a...))
}

func lightred(format string, a ...interface{}) string {
	return colored(bunt.LightSalmon, render(format, a...))
}

func dimgray(format string, a ...interface{}) string {
	return colored(bunt.DimGray, render(format, a...))
}

func bold(format string, a ...interface{}) string {
	return bunt.Style(
		fmt.Sprintf(format, a...),
		bunt.EachLine(),
		bunt.Bold(),
	)
}

func italic(format string, a ...interface{}) string {
	return bunt.Style(
		render(format, a...),
		bunt.EachLine(),
		bunt.Italic(),
	)
}

func pathToString(path *ytbx.Path, useGoPatchPaths bool, showPathRoot bool) string {
	var result string

	if useGoPatchPaths {
		result = styledGoPatchPath(path)

	} else {
		result = styledDotStylePath(path)
	}

	if path != nil && showPathRoot {
		result += bunt.Sprintf("  LightSteelBlue{(%s)}", path.RootDescription())
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
