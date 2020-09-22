package activity

import (
	"fmt"
	"testing"
)

func TestParam_String(t *testing.T) {

	tests := []struct {
		name  string
		param *Param
		want  string
	}{
		{"param", NewParam("key", "value"), "key=value"},
		{"param with nil value", NewParam("empty", nil), "empty"},
		{"param with empty name", NewParam("", 10), "10"},
		{"param with empty pair", NewParam("", nil), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.param.String()
			if got != tt.want {
				t.Errorf("NewParam.String() = '%v', want '%v'", got, tt.want)
			} else {
				fmt.Printf("NewParam.String() = '%v'\n", got)
			}
		})
	}
}
