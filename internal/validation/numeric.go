package validation

import "math"

func Clamp(value, low, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var total float64
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}
func Spread(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	low, high := values[0], values[0]
	for _, value := range values[1:] {
		low = math.Min(low, value)
		high = math.Max(high, value)
	}
	return high - low
}
func HueDistance(a, b float64) float64 {
	delta := math.Abs(a - b)
	if delta > 180 {
		delta = 360 - delta
	}
	return delta
}
func StableScore(hue, target, tolerance, ph, targetPH float64) float64 {
	return (Clamp(100-HueDistance(hue, target)/math.Max(tolerance, .1)*40, 0, 100) + Clamp(100-math.Abs(ph-targetPH)*20, 0, 100)) / 2
}
