package salute

import (
	"go-cases/salute/hello"
	"go-cases/salute/say"
)

// Print salute
func Print() string {
	return say.Print() + " " + hello.Print()
}
