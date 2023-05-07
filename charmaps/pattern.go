package charmaps

import (
	"errors"
	"strings"
)

type Pattern struct {
	CodePoint string
	Pixel     string
	Mask      uint64
	SetNum    int
}

func NewPattern(codePoint string, pixel string) (Pattern, error) {
	pat := strings.ReplaceAll(pixel, "\n", "")

	if len(pat) != 8*8 {
		return Pattern{}, errors.New("pattern is not 8x8 long")
	}

	// As we have 8x8 (=64) pixels we can use an uint64 and set bits to 1 where
	// the pixel is not a space. That way we can quickly check if a pixel is set.
	mask := uint64(0)
	for i := range pat {
		if pat[i] == []byte("X")[0] {
			mask = mask | (1 << (64 - i))
		}
	}

	return Pattern{
		CodePoint: codePoint,
		Pixel:     strings.ReplaceAll(pixel, "\n", ""),
		Mask:      mask,
		SetNum:    len(pat) - strings.Count(pat, " "),
	}, nil
}

func MustNewPattern(codePoint string, pixel string) Pattern {
	pat, err := NewPattern(codePoint, pixel)
	if err != nil {
		panic(err)
	}
	return pat
}

func (p Pattern) IsSet2(x int, y int) bool {
	return p.IsSet(x + y*8)
}

func (p Pattern) IsSet(i int) bool {
	return p.Mask&(1<<(64-i)) >= 1
}
