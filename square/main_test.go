package main

import (
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {

	tests := []struct {
		name string
		src  []int
		want []int
	}{
		{"-+", []int{-3, 2, 4}, []int{4, 9, 16}},
		{"nil", []int{}, []int{}},
		{"-", []int{-3, -2, -1}, []int{1, 4, 9}},
		{"+", []int{1, 2, 3}, []int{1, 4, 9}},
		{"0", []int{0}, []int{0}},
		{"-4 4", []int{-4, 4}, []int{16, 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Square(tt.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Square() = %v, want %v", got, tt.want)
			}
		})
	}
}
