package main

import (
	"fmt"
	"log"
	"koneko/source/hypr"
)

func main() {
	x, y, err := hypr.GetCursorPos()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d, %d", x, y)
}
