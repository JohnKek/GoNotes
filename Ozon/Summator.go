package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n int
	var in *bufio.Reader
	var out *bufio.Writer
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n)
	var a, b int
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a, &b)
		fmt.Fprintln(out, a+b)
	}

}
