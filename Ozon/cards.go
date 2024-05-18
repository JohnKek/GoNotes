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
	var n, m int
	fmt.Fscan(in, &n, &m)

	friends_card := make([]Card, n)
	var str string
	for j := 0; j < n; j++ {
		fmt.Fscan(in, &str)
		number, _ := strconv.Atoi(str)
		friends_card = append(friends_card, Card{j, number})
	}
	friends_card = friends_card[n:]
	sort.Sort(ByMax(friends_card))
	result := make([]int, n)
	flag := true
	for _, card := range friends_card {
		if card.max < m {
			result[card.index] = m
			m--
		} else {
			flag = false
			break
		}
	}
	if flag {
		for _, card := range result {
			fmt.Print(card) // Выводим places[key] без перехода на новую строку
			fmt.Print(" ")
		}
	} else {
		fmt.Print(-1)
	}
}

type Card struct {
	index int
	max   int
}

type ByMax []Card

func (a ByMax) Len() int           { return len(a) }
func (a ByMax) Less(i, j int) bool { return a[i].max > a[j].max }
func (a ByMax) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
