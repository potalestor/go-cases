package main

import "fmt"

/*
Пример
[-3, 2, 4] -> [4, 9, 16]

Сложность по времени не хуже O(N * log(N))
|xi| < 10000
*/

func middle(a []int) int {
	for i := range a {
		if a[i] >= 0 {
			return i
		}
	}
	return -1
}

// Square ...
func Square(a []int) []int {
	if len(a) == 0 {
		return a
	}
	result := []int{}

	m := middle(a)
	n := []int{}
	p := []int{}

	if m == -1 {
		n = a
		for i := len(a) - 1; i >= 0; i-- {
			result = append(result, n[i]*n[i])
		}
		return result

	} else if m == 0 {
		p = a
		for i := 0; i < len(a); i++ {
			result = append(result, p[i]*p[i])
		}
		return result

	} else {
		n := a[:m]
		p = a[m:]

		i := 0
		j := len(n) - 1
		for {
			if (j >= 0) && (i < len(p)) {
				if n[j]*n[j] <= p[i]*p[i] {
					result = append(result, n[j]*n[j])
					j--
					continue
				}
				result = append(result, p[i]*p[i])
				i++
				continue
			}
			if (j < 0) && (i < len(p)) {
				result = append(result, p[i]*p[i])
				i++
				continue
			}
			if (j >= 0) && (i >= len(p)) {
				result = append(result, n[j]*n[j])
				j--
				continue
			}

			break

		}
		return result

	}

}

func main() {
	fmt.Println(Square([]int{-3, 2, 4}))
	fmt.Println(Square([]int{}))
	fmt.Println(Square([]int{-3, -2, -1}))
	fmt.Println(Square([]int{1, 2, 3}))
	fmt.Println(Square([]int{-3, -3}))
	fmt.Println(Square([]int{2, 2}))
	fmt.Println(Square([]int{1}))
	fmt.Println(Square([]int{-1}))
	fmt.Println(Square([]int{0}))

}
