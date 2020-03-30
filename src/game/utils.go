package game

import (
	"fmt"
	"math/rand"
)

func defaultFuzzy(v float64) string {
	if v < 0 || v > 1 {
		panic(fmt.Sprintf("invalid argument to defaultFuzzy: %f", v))
	}

	switch {
	case v >= .9:
		return "++"
	case v >= .7:
		return "+"
	case v >= .3:
		return "o"
	case v >= .1:
		return "-"
	default:
		return "--"
	}
}

func randomBool(t float64) bool {
	return rand.Float64() < t
}

func randomFloat(min, max float64) float64 {
	if min < 0 || max > 1 || min > max {
		panic(fmt.Sprintf("invalid arguments: %f, %f", min, max))
	}

	return rand.Float64()*(max-min) + min
}
