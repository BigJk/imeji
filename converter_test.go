package imeji

import (
	"bytes"
	"github.com/anthonynsimon/bild/transform"
	"image/png"
	"os"
	"testing"
)

func BenchmarkConversion(b *testing.B) {
	data, err := os.ReadFile("./test/image.png")
	if err != nil {
		b.Fatal(err)
	}

	img, err := png.Decode(bytes.NewBuffer(data))
	if err != nil {
		b.Fatal(err)
	}

	img = transform.Resize(img, 100*8, 25*8, transform.NearestNeighbor)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ImageString(img, WithMaxRoutines(1))
		if err != nil {
			b.Fatal(err)
		}
	}
}
