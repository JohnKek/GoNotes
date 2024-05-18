package main

import (
	"fmt"
)

func increaseSequence(sequence []int) []int {
	n := len(sequence)
	for i := n - 1; i > 0; i-- {
		if sequence[i] > sequence[i-1] {
			for j := i; j < n; j++ {
				if sequence[j] > sequence[i-1] && (j == n-1 || sequence[j+1] <= sequence[i-1]) {
					sequence[i-1], sequence[j] = sequence[j], sequence[i-1]
					break
				}
			}
			for j, k := i, n-1; j < k; j, k = j+1, k-1 {
				sequence[j], sequence[k] = sequence[k], sequence[j]
			}
			return sequence
		}
	}
	return nil
}

func main() {
	sequence := []int{2, 1, 2, 3, 5}
	newSequence := increaseSequence(sequence)
	fmt.Println(newSequence)
}
