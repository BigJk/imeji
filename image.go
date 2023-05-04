package imeji

import (
	"github.com/BigJk/imeji/charmaps"
	"image"
	"image/color"
)

// pixelChunk Fetches a 8x8 image chunk from a image.
func pixelChunk(img image.Image, bounds image.Rectangle, cx, cy int) []color.Color {
	pixel := make([]color.Color, 8*8)
	for yi := 0; yi < 8; yi++ {
		for xi := 0; xi < 8; xi++ {
			pixel[xi+yi*8] = img.At(bounds.Min.X+cx*8+xi, bounds.Min.Y+cy*8+yi)
		}
	}
	return pixel
}

// patternError calculates how matching a character with foreground and background color is to the actual 8x8 pixels.
func patternError(pixel []color.Color, fg color.Color, bg color.Color, pattern *charmaps.Pattern) float64 {
	errSum := 0.0

	for i := 0; i < len(pattern.Pixel); i++ {
		if pattern.IsSet(i) {
			errSum += diff(pixel[i], fg)
		} else {
			errSum += diff(pixel[i], bg)
		}
	}

	return errSum
}
