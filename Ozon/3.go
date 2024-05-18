package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var u int
	fmt.Fscan(in, &u)
	var str string
	for i := 0; i < u; i++ {
		fmt.Fscan(in, &str)
		runes := []rune(str)
		flag := true
		etalon := runes[0]
		for index, _ := range runes {
			if index != len(runes)-1 && index != 0 && runes[index] != etalon {
				if runes[index-1] != runes[index+1] {
					flag = false
				}
			} else {
				if runes[index] != runes[0] {
					flag = false
				}
			}

		}
		if flag {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
