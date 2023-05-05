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

func sq(v uint32) uint32 {
	return v * v
}

func diff(a color.Color, b color.Color) float64 {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	return float64(sq(ar-br) + sq(ag-bg) + sq(ab-bb))
}

func diffSlow(a color.Color, b color.Color) float64 {
	ac, _ := colorful.MakeColor(a)
	bc, _ := colorful.MakeColor(b)
	return ac.DistanceRgb(bc)
}

func diffSlowest(a color.Color, b color.Color) float64 {
	ac, _ := colorful.MakeColor(a)
	bc, _ := colorful.MakeColor(b)
	return ac.DistanceLab(bc)
}
