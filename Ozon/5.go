package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	var in *bufio.Reader = bufio.NewReader(os.Stdin)
	var out *bufio.Writer = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var lines, n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lines)
		var str strings.Builder
		for i := 0; i < lines+1; i++ {
			byteLine, err := in.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			str.Write(bytes.TrimSpace([]byte(byteLine)))
		}
		var MyMap map[string]interface{}
		fmt.Println(str.String())
		err := json.Unmarshal([]byte("{"+str.String()+"}"), &MyMap)

		if err != nil {
			err = nil
		}
		fmt.Fprintln(out, MyMap)
	}
}
