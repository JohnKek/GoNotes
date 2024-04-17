package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var str string

	var in *bufio.Reader
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &str)

	fmt.Fscan(in, &n)

	for i := 0; i < n; i++ {
		var a, b int
		var s string
		fmt.Fscan(in, &a, &b, &s)
		str = str[:a-1] + s + str[b:]
	}

	// Выводим итоговую строку
	fmt.Fprintln(out, str)
}
