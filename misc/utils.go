package misc

func Sum(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	} else {
		return a[0] + Sum(a[1:])
	}
}
