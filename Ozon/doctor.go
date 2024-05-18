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
	var n, m, num int
	fmt.Fscan(in, &num)
	for i := 0; i < num; i++ {
		fmt.Fscan(in, &n, &m)
		slots := make(map[int]bool, m)
		for i := 1; i < m+1; i++ {
			slots[i] = true
		}
		patients := make([]int, n)
		var str string
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &str)
			number, _ := strconv.Atoi(str)
			patients = append(patients, number)
			delete(slots, number)
		}
		patients = patients[n:]
		result_string := ""
		windows := make(map[int]bool, m)
		deleteKey := 0
		for _, pat := range patients {
			if _, ok := windows[pat]; ok {
				if len(patients) != 0 {
					for key := range slots {
						deleteKey = key
						if key > pat {
							result_string += "+"
						} else {
							result_string += "-"
						}
						break
					}
					delete(slots, deleteKey)

				} else {
					result_string = "x"
					break
				}

			} else {
				windows[pat] = true
				result_string += "0"
			}
		}
		fmt.Println(result_string)

	}
}
