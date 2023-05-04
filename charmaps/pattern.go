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

	mask := uint64(0)
	for i := range pat {
		if pat[i] != []byte(" ")[0] {
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
	return string(p.Pixel[x+y*8]) != " "
}

func (p Pattern) IsSet(i int) bool {
	return string(p.Pixel[i]) != " "
}
