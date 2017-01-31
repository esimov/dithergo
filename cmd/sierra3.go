package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Sierra-3",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0 },
				[]float32{ 2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0 },
				[]float32{ 0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}