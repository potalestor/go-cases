package addition

// // C - код
// #include <stdio.h>
// void add(int a, int b) {
//     printf("%d\n", a + b);
// }
import "C"

// Go-код
import "fmt"

// Add 2 numbers
func Add(a, b int) {
	fmt.Printf("%d + %d = ", a, b)
	C.add(C.int(a), C.int(b))
}
