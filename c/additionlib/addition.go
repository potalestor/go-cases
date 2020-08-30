package additionlib

// // C - код
// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L. -laddition
// #include "addition.h"
import "C"

// Go-код
import "fmt"

// Add 2 numbers
func Add(a, b int) {
	fmt.Printf("%d + %d = ", a, b)
	C.add(C.int(a), C.int(b))
}
