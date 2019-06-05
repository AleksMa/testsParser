package main

import (
	"fmt"
	"time"
)

func kek() {
	//layout1 := "1136214245"
	//layout2 := "Monday, 2-Jan-06 15:04:05 UTC"


	/*i, _ := strconv.Atoi("946684810")
	T1 := time.Date(1970, 1, 1, 0, 0, i, 0, time.UTC)*/

	T1 := time.Unix(946684810, 0)


	/*if err1 != nil {
		fmt.Println(err1)
	}*/
	T2, err2 := time.Parse(time.RFC850, "Saturday, 01-Jan-00 00:00:20 UTC")
	if err2 != nil {
		fmt.Println(err2)
	}
	T3, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:10+00:00")
	//sliceVar2, _ := json.Marshal(sliceVar1)
	//fmt.Println(sliceVar1)
	fmt.Println(T1)
	fmt.Println(T2)
	fmt.Println(T3)
	fmt.Println(T1.Equal(T3))
}
