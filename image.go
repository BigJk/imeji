package imeji

import (
	"github.com/BigJk/imeji/charmaps"
	"image"
	"image/color"
	"math"
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
func patternError(pixel []color.Color, fg color.Color, bg color.Color, pattern *charmaps.Pattern, max float64) float64 {
	// TODO: Improve this with SIMD.
	//
	// We can load the pixel and selected pattern as 2 float slice and use SIMD to do the diffing between the pixel
	// data in chunks. This would greatly improve the performance, as we can process more pixel with less CPU cycles.
	// Different SIMD approaches for different Arch's would be needed. NEON for arm64 and AVX512, AVX, SSE etc. for amd64.
	//
	// See: https://gorse.io/posts/avx512-in-golang.html

	errSum := 0.0

	for i := 0; i < len(pattern.Pixel); i++ {
		// If the pattern position is set we know this pixel needs to be compared with the foreground.
		if pattern.IsSet(i) {
			errSum += diff(pixel[i], fg)
		} else {
			errSum += diff(pixel[i], bg)
		}

		// If the error sum is already bigger than the max error we can stop early.
		if errSum > max {
			return math.MaxFloat64
		}
	}

	return errSum
}
