package main

func Sum(xs [8]int) int {
	total := 0
	for _, value := range xs {
		total += value
	}
	return total

}

func average(xs [8]int) float64 {
	total := 0
	for _, value := range xs {
		total += value
	}
	return float64(total) / float64(len(xs))
}

func averageFloat(ys [8]float64) float64 {
	total := 0.0
	for _, value := range ys {
		total += value
	}
	return float64(total) / float64(len(ys))
}

func reverse(xs [8]int) [8]int {
	var total [8]int
	a := 0
	for i := len(xs) - 1; i >= 0; i-- {
		total[a] = xs[i]
		a++
	}
	return total
}
