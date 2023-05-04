package imeji

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
)

func darkest(colors ...color.Color) color.Color {
	var darkest color.Color

	darkestValue := 101.0
	for i := range colors {
		col, _ := colorful.MakeColor(colors[i])
		_, _, l := col.Hsl()
		if l < darkestValue {
			darkestValue = l
			darkest = colors[i]
		}
	}

	return darkest
}

func lightest(colors ...color.Color) color.Color {
	var lightest color.Color

	lightestValue := -1.0
	for i := range colors {
		col, _ := colorful.MakeColor(colors[i])
		_, _, l := col.Hsl()
		if l > lightestValue {
			lightestValue = l
			lightest = colors[i]
		}
	}

	return lightest
}

func diff(a color.Color, b color.Color) float64 {
	colA, _ := colorful.MakeColor(a)
	colB, _ := colorful.MakeColor(b)

	return colA.DistanceRgb(colB)
}
