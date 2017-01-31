package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Stucki",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0 },
				[]float32{ 2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0 },
				[]float32{ 1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}