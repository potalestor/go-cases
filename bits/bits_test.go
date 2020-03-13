package bits

import (
	"fmt"
	"testing"
)

func TestSetBit(t *testing.T) {
	type args struct {
		n   int
		pos uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"8", args{0, 3}, 8},
		{"10", args{2, 3}, 10},
		{"15", args{15, 3}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetBit(tt.args.n, tt.args.pos)
			if got != tt.want {
				t.Errorf("SetBit(%v, %d) = %v, want %v", AsString(tt.args.n), tt.args.pos, AsString(got), AsString(tt.want))
			} else {
				fmt.Printf("SetBit(%v, %d) = %v\n", AsString(tt.args.n), tt.args.pos, AsString(got))
			}
		})
	}
}

func TestClearBit(t *testing.T) {
	type args struct {
		n   int
		pos uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"0", args{8, 3}, 0},
		{"10", args{10, 3}, 2},
		{"7", args{15, 3}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClearBit(tt.args.n, tt.args.pos)
			if got != tt.want {
				t.Errorf("ClearBit(%v, %d) = %v, want %v\n", AsString(tt.args.n), tt.args.pos, AsString(got), AsString(tt.want))
			} else {
				fmt.Printf("ClearBit(%v, %d) = %v\n", AsString(tt.args.n), tt.args.pos, AsString(got))
			}
		})
	}
}

func TestHasBit(t *testing.T) {
	type args struct {
		n   int
		pos uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"8t", args{8, 3}, true},
		{"10t", args{10, 3}, true},
		{"10f", args{10, 2}, false},
		{"10t", args{10, 1}, true},
		{"10f", args{10, 0}, false},
		{"7f", args{7, 3}, false},
		{"7t", args{7, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasBit(tt.args.n, tt.args.pos)
			if got != tt.want {
				t.Errorf("HasBit(%v, %d) = %v, want %v\n", AsString(tt.args.n), tt.args.pos, got, tt.want)
			} else {
				fmt.Printf("HasBit(%v, %d) = %v\n", AsString(tt.args.n), tt.args.pos, got)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		hi    int
		lo    int
		shift uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"10", args{2, 2, 2}, 10},
		{"6", args{3, 2, 1}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Join(tt.args.hi, tt.args.lo, tt.args.shift)
			if got != tt.want {
				t.Errorf("Join(%v, %v, %d) = %v, want %v\n", AsString(tt.args.hi), AsString(tt.args.lo), tt.args.shift, AsString(got), AsString(tt.want))
			} else {
				fmt.Printf("Join(%v, %v, %d) = %v\n", AsString(tt.args.hi), AsString(tt.args.lo), tt.args.shift, AsString(got))
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		n     int
		shift uint
	}
	tests := []struct {
		name   string
		args   args
		wantHi int
		wantLo int
	}{
		{"10", args{10, 2}, 2, 2},
		{"13", args{13, 2}, 3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHi, gotLo := Split(tt.args.n, tt.args.shift)
			if gotHi != tt.wantHi {
				t.Errorf("Split(%v, %d) gotHi = %v, want %v\n", AsString(tt.args.n), tt.args.shift, AsString(gotHi), AsString(tt.wantHi))
			} else if gotLo != tt.wantLo {
				t.Errorf("Split(%v, %d) gotLo = %v, want %v\n", AsString(tt.args.n), tt.args.shift, AsString(gotLo), AsString(tt.wantLo))
			} else {
				fmt.Printf("Split(%v, %d) hi = %v lo = %v\n", AsString(tt.args.n), tt.args.shift, AsString(gotHi), AsString(gotLo))
			}
		})
	}
}
