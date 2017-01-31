package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Atkinson",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 1.0 / 8.0, 1.0 / 8.0 },
				[]float32{ 1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0.0 },
				[]float32{ 0.0, 1.0 / 8.0, 0.0, 0.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}