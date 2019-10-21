package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var files []string
	var buf []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var counter int
			for _, f := range strings.Split(path, "/") {
				if stringInSlice(f, buf) {
					counter += 4
				} else if strings.Index(f, ".") != -1 {
					space := strings.Repeat(" ", counter)
					files = append(files, space+"├───"+f)
					counter = 0
				} else {
					space := strings.Repeat(" ", counter)
					files = append(files, space+"└───"+f)
					buf = append(buf, f)
					counter += 4
				}
			}
			counter = 0
		}
		return nil
	})
	for _, path := range files {
		fmt.Println(path)
	}
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
