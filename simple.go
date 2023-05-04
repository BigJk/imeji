package imeji

import (
	"bytes"
	"github.com/muesli/termenv"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

var codes = map[int]string{
	0b0000: " ",
	0b0001: "▗",
	0b0010: "▖",
	0b0011: "▄",
	0b0100: "▝",
	0b0101: "▐",
	0b0110: "▞",
	0b0111: "▟",
	0b1000: "▘",
	0b1001: "▚",
	0b1010: "▌",
	0b1011: "▙",
	0b1100: "▀",
	0b1101: "▜",
	0b1110: "▛",
	0b1111: "█",
}

func FileSimple(out io.Writer, path string, options ...Option) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	return ImageSimple(out, img, options...)
}

func ImageSimple(out io.Writer, img image.Image, options ...Option) error {
	optionData := &OptionData{
		fontRatio: 0,
		img:       img,
		out:       termenv.NewOutput(out, termenv.WithProfile(termenv.EnvColorProfile())),
		buf:       out,
	}

	for i := range options {
		options[i](optionData)
	}

	bounds := optionData.img.Bounds().Max

	chunkX := bounds.X / 2
	chunkY := bounds.Y / 2

	for y := 0; y < chunkY; y++ {
		for x := 0; x < chunkX; x++ {
			tl := optionData.img.At(x*2, y*2)
			tr := optionData.img.At(x*2+1, y*2)
			bl := optionData.img.At(x*2, y*2+1)
			br := optionData.img.At(x*2+1, y*2+1)

			dark := darkest(tl, tr, bl, br)
			light := lightest(tl, tr, bl, br)

			mask := 0

			if diff(tl, light) < diff(tl, dark) {
				mask = mask | 0b1000
			}

			if diff(tr, light) < diff(tr, dark) {
				mask = mask | 0b0100
			}

			if diff(bl, light) < diff(bl, dark) {
				mask = mask | 0b0010
			}

			if diff(br, light) < diff(br, dark) {
				mask = mask | 0b0001
			}

			if _, err := optionData.out.Write([]byte(optionData.out.String(codes[mask]).Background(optionData.out.FromColor(dark)).Foreground(optionData.out.FromColor(light)).String())); err != nil {
				return err
			}
		}

		if _, err := out.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

func FileStringSimple(path string, options ...Option) (string, error) {
	buf := &bytes.Buffer{}
	err := FileSimple(buf, path, options...)
	return buf.String(), err
}

func ImageStringSimple(img image.Image, options ...Option) (string, error) {
	buf := &bytes.Buffer{}
	err := ImageSimple(buf, img, options...)
	return buf.String(), err
}
