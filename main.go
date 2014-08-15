package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Generate greyscale palette from single float f.
func fromPalette(f float64) (col color.NRGBA, err error) {
	if f < 0 || f > 1 {
		return color.NRGBA{0, 0, 0, 0}, errors.New("Number must be between zero and one")
	}

	// Scale f from 0 to 255 with zero f being white
	scaled := uint8(255 - (f * 255))
	return color.NRGBA{scaled, scaled, scaled, 255}, nil
}

// Computes whether complex number c is in the mandlebrot set.
func inSet(c complex128, maxIters int) (inSet bool, iters int) {
	z := complex128(0)
	for iters = 0; iters < maxIters; iters++ {
		// Here is the magic!
		z = z*z + c
		// Break if our number exceeds escape conditions
		if cmplx.Abs(z) > 60 {
			inSet = false
			return
		}
	}
	inSet = true
	return
}

var SIZE int = 1024

func main() {
	bounds := image.Rect(0, 0, SIZE, SIZE)
	m := image.NewNRGBA(bounds)

	// Smaller is bigger zoom
	zoom := 2.0
	offsetX, offsetY := 1.5, 1.5
	maxIters := 30

	// Make our fractal
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			// Generate our complex number, using offsets and scaling (zoom)
			c := complex(float64(x)/(float64(bounds.Max.X)/zoom)-offsetX,
				float64(y)/(float64(bounds.Max.Y)/zoom)-offsetY)

			_, iters := inSet(c, maxIters)
			f := float64(iters) / float64(maxIters)
			col, err := fromPalette(f)
			check(err)
			// Set pixel in our image
			m.Set(x, y, col)
		}
	}

	f, err := os.Create("out.png")
	check(err)

	defer f.Close()

	// Write the file to output
	err = png.Encode(f, m)
	check(err)

	f.Sync()
}
