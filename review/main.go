package main

import "fmt"

func Sum(xs [8]int) int {
	total := 0
	for _, value := range xs {
		total += value
	}
	return total

}

func Average(xs [8]int) float64 {
	total := 0
	for _, value := range xs {
		total += value
	}
	return float64(total) / float64(len(xs))
}

func AverageFloat(ys [8]float64) float64 {
	total := 0.0
	for _, value := range ys {
		total += value
	}
	return float64(total) / float64(len(ys))
}

func Reverse(xs [8]int) [8]int {
	var total [8]int
	a := 0
	for i := len(xs) - 1; i >= 0; i-- {
		total[a] = xs[i]
		a++
	}
	return total
}

// Примеры использования функций:

func main() {
	xs := [8]int{1, 2, 3, 4, 5, 6, 7, 8}

	fmt.Println(Sum(xs)) // Вывод: 36

	fmt.Println(Average(xs)) // Вывод: 4.5

	ys := [8]float64{1, 2, 3, 4, 5, 6, 7, 8}

	fmt.Println(AverageFloat(ys)) // 4.5

	fmt.Println(Reverse(xs)) // Вывод: [8 7 6 5 4 3 2 1]
}
