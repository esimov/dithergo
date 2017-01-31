package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Sierra-2",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 0.0, 4.0 / 16.0, 3.0 / 16.0 },
				[]float32{ 1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0 },
				[]float32{ 0.0, 0.0, 0.0, 0.0, 0.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}