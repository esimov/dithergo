package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"Sierra-Lite",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 2.0 / 4.0 },
				[]float32{ 1.0 / 4.0, 1.0 / 4.0, 0.0 },
				[]float32{ 0.0, 0.0, 0.0 },
			},
		},
	},
}

func main() {
	dither.Process(ditherers)
}