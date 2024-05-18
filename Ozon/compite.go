package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		runners := make([]Result, num)
		for j := 0; j < num; j++ {
			fmt.Fscan(in, &str)
			number, _ := strconv.Atoi(str)
			runners = append(runners, Result{j, number})
		}
		runners = runners[num:]
		sort.Sort(ByTime(runners))
		place := 1
		count := 0
		pr_data := runners[0].time
		places := make(map[int]int, len(runners))
		for _, runner := range runners {
			if runner.time == pr_data || runner.time == pr_data+1 {
				places[runner.index] = place
				count++
				pr_data = runner.time
			} else {
				place = place + count
				places[runner.index] = place
				count = 1
				pr_data = runner.time
			}

		}
		for i := 0; i < num; i++ {
			fmt.Print(places[i]) // Выводим places[key] без перехода на новую строку
			fmt.Print(" ")       // Добавляем пробел после каждого элемента
		}
		fmt.Println()

	}

}

type Result struct {
	index int
	time  int
}

type ByTime []Result

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Less(i, j int) bool { return a[i].time < a[j].time }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
