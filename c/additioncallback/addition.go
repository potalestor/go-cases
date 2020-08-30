package additioncallback

/*
// C - код
#cgo LDFLAGS: -lpthread
#include <stdio.h>
extern void adds();
*/
import "C"

// Go-код
import "fmt"

//export Add2Numbers
func Add2Numbers(a C.int, b C.int) {
	fmt.Println(a + b)
}

// Add 2 numbers
func Add(a, b int) {
	C.adds()
}
