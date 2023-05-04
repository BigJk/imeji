package imeji

import (
	"github.com/BigJk/imeji/charmaps"
	"github.com/anthonynsimon/bild/transform"
	"github.com/muesli/termenv"
	"image"
	"io"
)

type OptionData struct {
	pattern    []charmaps.Pattern
	routines   int
	colorPairs int
	fontRatio  float64
	img        image.Image
	out        *termenv.Output
	buf        io.Writer
}

type Option func(data *OptionData)

func WithResize(width int, height int) Option {
	return func(data *OptionData) {
		data.img = transform.Resize(data.img, width*8, height*8, transform.NearestNeighbor)
	}
}

// WithMaxWidth specifies a max width in cells for the output. This will keep the aspect ratio of the input picture and
// scale based on specified font ratio.
func WithMaxWidth(width int) Option {
	return func(data *OptionData) {
		rect := (data.img).Bounds()
		if rect.Max.X > width {
			aspectRatio := float64(rect.Max.Y) / float64(rect.Max.X)
			ratio := float64(width) / float64(rect.Max.X) * data.fontRatio

			data.img = transform.Resize(data.img, width*8, int(float64(rect.Max.Y)*ratio*aspectRatio)*8, transform.Linear)
		}
	}
}

// WithCrop crops a part of the image. This uses the normal image coordinates.
func WithCrop(x int, y int, width int, height int) Option {
	return func(data *OptionData) {
		data.img = transform.Crop(data.img, image.Rect(x, y, x+width, y+height))
	}
}

// WithTrueColor enables true color support. This is important if you don't output to os.Stdout and no terminal
// detection can be done or if you want to force a color mode.
func WithTrueColor() Option {
	return func(data *OptionData) {
		data.out = termenv.NewOutput(data.buf, termenv.WithProfile(termenv.TrueColor))
	}
}

// WithANSI enables basic ansi color support. This is important if you don't output to os.Stdout and no terminal
// detection can be done or if you want to force a color mode.
func WithANSI() Option {
	return func(data *OptionData) {
		data.out = termenv.NewOutput(data.buf, termenv.WithProfile(termenv.ANSI))
	}
}

// WithANSI256 enables 256 color support. This is important if you don't output to os.Stdout and no terminal
// detection can be done or if you want to force a color mode.
func WithANSI256() Option {
	return func(data *OptionData) {
		data.out = termenv.NewOutput(data.buf, termenv.WithProfile(termenv.ANSI256))
	}
}

// WithMaxRoutines specifies how many go routines are allowed to be spawned for calculating the image.
func WithMaxRoutines(routines int) Option {
	return func(data *OptionData) {
		data.routines = routines
	}
}

// WithColorPairMax specifies how many color pair possibilities the algorithm will try per 8x8 pixel chunk. Lower value
// results in better performance but colors and selected symbols might be suboptimal. Values between 1 and 12 are sensible.
// Default is 6.
func WithColorPairMax(pairs int) Option {
	return func(data *OptionData) {
		data.colorPairs = pairs
	}
}

// WithFontScaling sets the vertical font scaling size, as most terminal fonts are taller than wider.
func WithFontScaling(scale float64) Option {
	return func(data *OptionData) {
		data.fontRatio = scale
	}
}

// WithPattern specifies the character patterns that are usable in the algorithms. More patterns decrease the performance.
func WithPattern(pattern ...[]charmaps.Pattern) Option {
	return func(data *OptionData) {
		data.pattern = charmaps.Combine(pattern...)
	}
}
