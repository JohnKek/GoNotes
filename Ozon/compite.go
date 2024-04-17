package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	var num int
	var str string
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &num)
		time := make([]int, num)
		place := make([]int, num)
		for j := 0; j < num; j++ {
			fmt.Fscan(in, &str)
			number, _ := strconv.Atoi(str)
			time[number] = number
		}
		for {

		}

	}

}
func Max(arr []int) int {
	max := arr[0]
	for index, value := range arr {
		if value > max {
			max = index
		}
	}
	return max
}
