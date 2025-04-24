package main

import (
	"fmt"
	"time"
)

func main() {
	i := 1738000440
	fmt.Println(time.Now())
	fmt.Println(NowToTime(int64(i)))
}

func NowToTime(unixTime int64) *time.Time {
	if unixTime == 0 {
		return nil
	}
	t := time.Unix(unixTime, 0).UTC()
	return &t
}
