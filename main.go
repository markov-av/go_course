package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	alreadySeen := make(map[string]bool)
	for in.Scan() {
		txt := in.Text()
		if _, found := alreadySeen[txt]; found {
			continue
		}
		alreadySeen[txt] = true
		fmt.Println(txt)
	}
}
