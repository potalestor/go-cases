// Package bits implements bit manipulation functions for the integer type
package bits

import "fmt"

// SetBit - set the bit at pos in the integer n
func SetBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// ClearBit - clear the bit at pos in n
func ClearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

// HasBit - check the bit at pos in n
func HasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

// Join - join 2 number at one
func Join(hi, lo int, shift uint) int {
	return (hi << shift) | lo
}

// Split - split the number to hi and lo parts using shift
func Split(n int, shift uint) (hi, lo int) {
	hi = n >> shift
	lo = n & hi
	return
}

// AsString - return int as string
func AsString(n int) string {
	return fmt.Sprintf("%d (%bb)", n, n)
}
