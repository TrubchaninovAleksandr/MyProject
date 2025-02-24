package main

import (
	"bytes"
	"os"
	"testing"
)

func TestSum(t *testing.T) {
	numbers := [8]int{5, 6, 7, 8, 9, 10, 11, 12}
	otvet := 68
	result := Sum(numbers)

	if result != otvet {
		t.Errorf("Тест не пройден Sum(%v) = %d; ожидается %d", numbers, otvet, result)
	}
}

func TestAverage(t *testing.T) {
	numbers := [8]int{5, 6, 7, 8, 9, 10, 11, 12}
	otvet := 8.5
	result := Average(numbers)

	if float64(result) != otvet {
		t.Errorf("Тест не пройден Sum(%v) = %v; ожидается %v", numbers, otvet, result)
	}
}

func TestAverageFloat(t *testing.T) {
	numbers := [8]float64{5, 6, 7, 8, 9, 10, 11, 12}
	otvet := 8.5
	result := AverageFloat(numbers)

	if result != otvet {
		t.Errorf("Тест не пройден Sum(%v) = %v; ожидается %v", numbers, otvet, result)
	}
}

func TestReverse(t *testing.T) {
	numbers := [8]int{5, 6, 7, 8, 9, 10, 11, 12}
	otvet := [8]int{12, 11, 10, 9, 8, 7, 6, 5}
	result := Reverse(numbers)
	if result != otvet {
		t.Errorf("Тест не пройден Sum(%v) = %#v; ожидается %#v", numbers, otvet, result)
	}
}
func TestMainFunc(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = old

	var stdout bytes.Buffer
	stdout.ReadFrom(r)

	expected := "36\n4.5\n4.5\n[8 7 6 5 4 3 2 1]\n"
	if stdout.String() != expected {
		t.Errorf("got %#v, want %#v", stdout.String(), expected)
	}
}
