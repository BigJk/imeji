package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/imeji"
	"github.com/BigJk/imeji/charmaps"
	"os"
	"strings"
)

func main() {
	input := flag.String("input", "", "input image path")
	symbols := flag.String("symbols", "blocks", "which symbole sets to use (blocks, blocks_simple, blocks_adv, ascii, misc)")
	size := flag.String("size", "", "size in terminal cells (e.g. 100x20)")
	mw := flag.Int("max-width", 0, "sets the max width of the output (in cells) and keeps the aspect ratio")
	sharpen := flag.Int("sharpen", 0, "sets the number of times the image will be sharpened")
	brightness := flag.Float64("brightness", 0, "sets the change of brightness")
	saturation := flag.Float64("saturation", 0, "sets the change of saturation")
	contrast := flag.Float64("contrast", 0, "sets the change of contrast")
	flipH := flag.Bool("flip-h", false, "flip image horizontal")
	flipV := flag.Bool("flip-v", false, "flip image vertical")
	scale := flag.Float64("font-scale", 0, "vertical font scaling value (default 0.8)")
	forceFullColor := flag.Bool("force-full-color", false, "forces full color output")
	help := flag.Bool("help", false, "print help")
	flag.Parse()

	if *help {
		fmt.Println("イメジ :: Images for the terminal ー by BigJk")
		fmt.Println(" _                 _ _ \n(_)_ __ ___   ___ (_|_)\n| | '_ ` _ \\ / _ \\| | |\n| | | | | | |  __/| | |\n|_|_| |_| |_|\\___|/ |_|\n                |__/   \n_________________________________________")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	if len(*input) == 0 {
		fmt.Println("error: specify input file")
		return
	}

	var options []imeji.Option

	if *scale > 0 {
		options = append(options, imeji.WithFontScaling(*scale))
	}

	if len(*symbols) > 0 {
		var sets [][]charmaps.Pattern
		setNames := strings.Split(*symbols, ",")

		for i := range setNames {
			if val, ok := charmaps.CharMaps[strings.ToLower(strings.TrimSpace(setNames[i]))]; ok {
				sets = append(sets, val)
			}
		}

		options = append(options, imeji.WithPattern(sets...))
	}

	if len(*size) > 0 {
		var width int
		var height int

		if n, err := fmt.Sscanf(*size, "%dx%d", &width, &height); n != 2 || err != nil {
			panic("size argument malformed")
		}

		options = append(options, imeji.WithResize(width, height))
	}

	if *mw > 0 {
		options = append(options, imeji.WithMaxWidth(*mw))
	}

	if *forceFullColor {
		options = append(options, imeji.WithTrueColor())
	}

	if *sharpen > 0 {
		options = append(options, imeji.WithSharpen(*sharpen))
	}

	if *contrast > 0 {
		options = append(options, imeji.WithContrast(*contrast))
	}

	if *brightness > 0 {
		options = append(options, imeji.WithBrightness(*brightness))
	}

	if *saturation > 0 {
		options = append(options, imeji.WithSaturation(*saturation))
	}

	options = append(options, imeji.WithFlip(*flipH, *flipV))

	err := imeji.File(os.Stdout, *input, options...)
	if err != nil {
		panic(err)
	}
}
