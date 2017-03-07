package main

import (
	"fmt"
	"math"
)

func fmtElapsedTime(seconds float32) string {
	units := []string{"Nanosecond", "Microsecond", "Millisecond", "Second", "Minute", "Hour", "Day"}
	divisors := []int{1, 1000, 1000, 1000, 60, 60, 24}
	v := math.Abs(float64(seconds)) * 1000.0 * 1000.0 * 1000.0
	lastU := units[0]
	for i := range units {
		if v < float64(divisors[i]) {
			break
		}
		lastU = units[i]
		v /= float64(divisors[i])
	}
	plural := "s"
	if v == 1.0 {
		plural = ""
	}
	return fmt.Sprintf("%.2f %v%v", v, lastU, plural)
}
