package activity

import (
	"fmt"
	"testing"
)

func Test_messageBuilder_buildStructuredData(t *testing.T) {
	tests := []struct {
		name string
		data []interface{}
		want string
	}{
		{"nil", nil, ""},
		{"string", []interface{}{"string"}, "string"},
		{"primitive types", []interface{}{"string", 10, 40.8}, "string 10 40.8"},
		{"pair", []interface{}{NewParam("key", "value")}, "key=value"},
		{"pairs", []interface{}{NewParam("key", "value"), NewParam("int", 10)}, "key=value int=10"},
		{"pairs with nil value", []interface{}{NewParam("key", "value"), NewParam("empty", nil)}, "key=value empty"},
		{"pairs with empty key", []interface{}{NewParam("key", "value"), NewParam("", 10)}, "key=value 10"},
		{"pairs with empty pair", []interface{}{NewParam("key", "value"), NewParam("", nil)}, "key=value"},
		{"pairs and primitive types", []interface{}{NewParam("key", "value"), NewParam("int", 10), "string", 20}, "key=value int=10 string 20"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMessageBuilder().buildStructuredData(tt.data...).build()
			got := m.StructuredData
			if got != tt.want {
				t.Errorf("messageBuilder.buildStructuredData() = '%v', want '%v'", got, tt.want)
			} else {
				fmt.Printf("messageBuilder.buildStructuredData() = '%v'\n", got)
			}
		})
	}
}

func Test_messageBuilder_buildProcessID(t *testing.T) {
	m := newMessageBuilder().buildProcessID().build()
	if m.ProcessID == 0 {
		t.Errorf("messageBuilder.buildProcessID() = '%v', want '>0'", m.ProcessID)
	} else {
		fmt.Printf("messageBuilder.buildProcessID() = '%v'\n", m.ProcessID)
	}
}

func Test_messageBuilder_buildHostname(t *testing.T) {
	m := newMessageBuilder().buildHostname().build()
	if m.Hostname == "" {
		t.Errorf("messageBuilder.buildHostname() = '%v', want !=''", m.Hostname)
	} else {
		fmt.Printf("messageBuilder.buildHostname() = '%v'\n", m.Hostname)
	}
}
