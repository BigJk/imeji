# imeji

イメジ ー Images for the terminal

![demo](./github/screenshot.png)

[![GoReportCard](https://goreportcard.com/badge/github.com/BigJk/imeji)](https://goreportcard.com/report/github.com/BigJk/imeji) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/BigJk/imeji)

imeji is a lightweight alternative to the awesome [chafa](https://hpjansson.org/chafa/). It is written in go and can be easily embedded into tools. imeji takes a image as input and creates a sequence of characters and ansi color sequences resulting in a terminal printable images.

### Why not use chafa?

If you can install chafa it should be preferred, as it is more advanced, faster and just awesome! But if you want to include terminal image output in your go application and don't want to ship chafa as external dependency imeji might be worth a try.

# CLI

```
イメジ :: Images for the terminal ー by BigJk
 _                 _ _
(_)_ __ ___   ___ (_|_)
| | '_ ` _ \ / _ \| | |
| | | | | | |  __/| | |
|_|_| |_| |_|\___|/ |_|
                |__/
_________________________________________

  -font-scale float
    	vertical font scaling value (default: 0.8)
  -force-full-color
    	forces full color output
  -help
    	print help
  -input string
    	input image path
  -max-width int
    	sets the max width of the output (in cells) and keeps the aspect ratio
  -size string
    	size in terminal cells (e.g. 100x20)
  -symbols string
    	which symbole sets to use (blocks, blocks_simple, blocks_adv, ascii, misc) (default "blocks")
```

### Install ``imeji`` command

```
go install github.com/BigJk/imeji/cmd/imeji@latest
```

# Go Library

```
go get github.com/BigJk/imeji
```

### Example

```go
// Print directly to stdout and detect terminal capabilities:
imeji.File(os.Stdout, "./image.png", imeji.WithMaxWidth(100))

// Convert to string with full color support:
text, _ := imeji.FileString("./image.png", imeji.WithTrueColor())
fmt.Println(text)
```

# Technique

I didn't find a good written reference on the technique used by Chafa and other tools so here is a basic overview. Its important to know that in the terminal we are limited to a single foreground and background color per character. That means for each cell in the terminal we need to find the best character and foreground, background pair with the least "error" (difference) to the real picture.

## Pattern

The basic idea is that you can map a single character in terminal to 8x8 pixels of a real image. For each character that wants to be used in a terminal picture the pattern needs to be created. A pattern can easily be defined by an 8 line string. The pattern defines which pixels are set to the foreground color and which to the background color.

See ``/charmaps/blocks.go`` and you will quickly get the idea.

**Examples**: 

- Char: █ (Full Block -> All pixel set)

```
XXXXXXXX
XXXXXXXX
XXXXXXXX
XXXXXXXX
XXXXXXXX
XXXXXXXX
XXXXXXXX
XXXXXXXX
```

- Char: !

```
________
___X____
___X____
___X____
___X____
________
___X____
________
```

## Procedure

Now that we know how to define the pattern for a character we can convert a image to terminal printable characters.

1. Chunk the image into 8x8 pixel chunks
2. For each pixel chunk:
   1. Randomly select a few (foreground, background) pairs from the pixels in the chunk
   2. For all the pairs test all characters and calculate the error to the real pixels
   3. Return the pair and character with the least error
3. For each chunk print the selected character with the chosen foreground and background

# Further Work

- [ ] Better handling of alpha channel
- [ ] Use assembler with SIMD instructions to improve pixel to pattern diffing performance