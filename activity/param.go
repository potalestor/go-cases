package activity

import "fmt"

// Param is "name=value" structure
type Param struct {
	Name  string
	Value interface{}
}

// NewParam create "name=value" structure
func NewParam(name string, value interface{}) *Param {
	return &Param{
		Name:  name,
		Value: value,
	}
}

func (p *Param) String() string {
	if (p.Name == "") && (p.Value == nil) {
		return ""
	}
	if p.Name == "" {
		return fmt.Sprintf("%v", p.Value)
	}
	if p.Value == nil {
		return p.Name
	}
	return fmt.Sprintf("%s=%v", p.Name, p.Value)
}
