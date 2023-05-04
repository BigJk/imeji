package imeji

import (
	"bytes"
	"github.com/BigJk/imeji/charmaps"
	"github.com/muesli/termenv"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type cellResult struct {
	fg    color.Color
	bg    color.Color
	point string
}

type workReq struct {
	x int
	y int
}

// BasicSelection finds the best character, foreground and background pair matching the 8x8 pixel grid.
func BasicSelection(options *OptionData, pixel []color.Color) (string, color.Color, color.Color) {
	pairs := options.colorPairs
	fgs := make([]color.Color, 0, pairs)
	bgs := make([]color.Color, 0, pairs)

	// Generate unique foreground background combos. Try at most pairs * 2 times.
	// We need to have a try limit in case the pixel have < pairs unique colors.
pairGen:
	for i := 0; i < pairs*2; i++ {
		fg := pixel[rand.Intn(len(pixel))]
		bg := pixel[rand.Intn(len(pixel))]

		// Check if unique
		for j := range fgs {
			if diff(fg, fgs[j]) == 0 && diff(bg, bgs[j]) == 0 {
				continue pairGen
			}
		}

		fgs = append(fgs, fg)
		bgs = append(bgs, bg)

		if len(fgs) >= pairs {
			break
		}
	}

	// Search for the character by minimizing against the error between character and pixels
	best := 99999.0
	point := " "
	index := 0
	for i := range fgs {
		for j := range options.pattern {
			err := patternError(pixel, fgs[i], bgs[i], &options.pattern[j])
			if err < best {
				best = err
				point = options.pattern[j].CodePoint
				index = i
			}
		}
	}

	return point, fgs[index], bgs[index]
}

// File decodes the image from a file and builds a terminal printable image of it.
func File(out io.Writer, path string, options ...Option) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	return Image(out, img, options...)
}

// Image decodes the images and builds a terminal printable image of it.
func Image(out io.Writer, img image.Image, options ...Option) error {
	optionData := &OptionData{
		pattern:    charmaps.BlocksBasic,
		fontRatio:  0.8,
		colorPairs: 6,
		routines:   runtime.NumCPU(),
		img:        img,
		out:        termenv.NewOutput(out, termenv.WithProfile(termenv.EnvColorProfile())),
		buf:        out,
	}

	for i := range options {
		options[i](optionData)
	}

	bounds := optionData.img.Bounds()
	chunkX := bounds.Dx() / 8
	chunkY := bounds.Dy() / 8

	res := make([]cellResult, chunkX*chunkY)
	workChan := make(chan workReq, chunkX*chunkY)

	wg := &sync.WaitGroup{}
	wg.Add(optionData.routines + 1)

	// Generate work request for each pixel and send it to the channel.
	go func() {
		defer wg.Done()

		for y := 0; y < chunkY; y++ {
			for x := 0; x < chunkX; x++ {
				workChan <- workReq{
					x: x,
					y: y,
				}
			}
		}

		// Wait for work channel to be fully drained
		for len(workChan) > 0 {
			time.Sleep(time.Millisecond)
		}

		close(workChan)
	}()

	// Start all the go routines and let them generate each pixel chunk
	for i := 0; i < optionData.routines; i++ {
		go func() {
			defer wg.Done()

			for req := range workChan {
				pixel := pixelChunk(optionData.img, bounds, req.x, req.y)

				point, fg, bg := BasicSelection(optionData, pixel)

				res[req.y*chunkX+req.x].fg = fg
				res[req.y*chunkX+req.x].bg = bg
				res[req.y*chunkX+req.x].point = point
			}
		}()
	}

	// Wait for finish
	wg.Wait()

	// Collect results
	for i := range res {
		if _, err := out.Write([]byte(optionData.out.String(res[i].point).Foreground(optionData.out.FromColor(res[i].fg)).Background(optionData.out.FromColor(res[i].bg)).String())); err != nil {
			return err
		}

		if i != 0 && (i+1)%chunkX == 0 {
			if _, err := out.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}

	return nil
}

// ImageString decodes the images and builds a terminal printable image of it. Returns the ansi string.
func ImageString(img image.Image, options ...Option) (string, error) {
	buf := &bytes.Buffer{}
	if err := Image(buf, img, options...); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FileString decodes the image from a file and builds a terminal printable image of it. Returns the ansi string.
func FileString(path string, options ...Option) (string, error) {
	buf := &bytes.Buffer{}
	if err := File(buf, path, options...); err != nil {
		return "", err
	}
	return buf.String(), nil
}
