package main

import (
	"reflect"
	"testing"
)

func Test_sum(t *testing.T) {
	type args struct {
		xs [8]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sum",
			args: args{[8]int{1, 2, 3, 4, 5, 6, 7, 8}},
			want: 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.xs); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_average(t *testing.T) {
	type args struct {
		xs [8]int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "average",
			args: args{[8]int{1, 2, 3, 4, 5, 6, 7, 8}},
			want: 4.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := average(tt.args.xs); got != tt.want {
				t.Errorf("average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_averageFloat(t *testing.T) {
	type args struct {
		ys [8]float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{

		{
			name: "average_float",
			args: args{[8]float64{1, 2, 3, 4, 5, 6, 7, 8}},
			want: 4.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := averageFloat(tt.args.ys); got != tt.want {
				t.Errorf("averageFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverse(t *testing.T) {
	type args struct {
		xs [8]int
	}
	tests := []struct {
		name string
		args args
		want [8]int
	}{
		{
			name: "revers",
			args: args{[8]int{1, 2, 3, 4, 5, 6, 7, 8}},
			want: [8]int{8, 7, 6, 5, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.args.xs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
