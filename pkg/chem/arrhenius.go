package chem

import "math"

const gasConstant = 8.314462618

type ArrheniusPoint struct {
	Temperature  float64 `json:"temperature"`
	RateConstant float64 `json:"rateConstant"`
}

func SweepArrhenius(A, Ea, Tmin, Tmax float64, points int) []ArrheniusPoint {
	if points < 2 {
		points = 2
	}
	step := (Tmax - Tmin) / float64(points-1)
	out := make([]ArrheniusPoint, 0, points)
	for i := 0; i < points; i++ {
		t := Tmin + float64(i)*step
		k := A * math.Exp(-Ea/(gasConstant*t))
		out = append(out, ArrheniusPoint{Temperature: t, RateConstant: k})
	}
	return out
}
