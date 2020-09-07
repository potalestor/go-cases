package say

import "testing"

func TestPrint(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Say", "Say"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Print(); got != tt.want {
				t.Errorf("Print() = %v, want %v", got, tt.want)
			}
		})
	}
}
