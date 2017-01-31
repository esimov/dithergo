package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Burkes",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0 },
				[]float32{ 2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0 },
				[]float32{ 0.0, 0.0, 0.0, 0.0, 0.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}