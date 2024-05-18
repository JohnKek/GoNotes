package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var u, n int
	fmt.Fscan(in, &u, &n)
	user := make(map[int][]int, u)
	var global_messages []int
	massege_number := 1
	var a, b int
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a, &b)
		switch a {
		case 1:
			switch b {
			case 0:
				global_messages = append(global_messages, massege_number)
				massege_number++
			default:
				user[b] = append(user[b], massege_number)
				massege_number++
			}
		case 2:
			local := 0
			global := 0
			if len(user[b]) != 0 {
				local = user[b][len(user[b])-1]
				if len(global_messages) != 0 {
					global = global_messages[len(global_messages)-1]

				}
				fmt.Fprintln(out, int(math.Max(float64(global), float64(local))))

			}
		}

	}
}
