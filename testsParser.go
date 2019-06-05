package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("1.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(data))
}
