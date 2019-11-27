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

func GetFileSize1(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var space string
	var words []string
	var byteSize int64
	var fileSize string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			byteSize, _ = GetFileSize1(path)
			if byteSize == 0 {
				fileSize = "(empty)"
			} else {
				fileSize = "(" + fmt.Sprintf("%v", byteSize) + "b)"
			}
			for idx, f := range strings.Split(path+fileSize, "/") {
				//if idx == 0 {
				//	continue
				//}
				space = strings.Repeat("|    ", idx)
				if contains(words, space+"├───"+f) {
					continue
				} else {
					words = append(words, space+"├───"+f)
				}
			}
		}
		return nil
	})
	var prevLen = 0
	var prevStr string
	for idx, str := range words {
		result := strings.Split(str, "├───")
		if prevLen > len(result[0]) {
			words[idx-1] = strings.Replace(prevStr, "├───", "└───", -1)
			prevLen = 0
		} else {
			prevStr = str
			prevLen = len(result[0])
		}
	}
	words[len(words)-1] = strings.Replace(words[len(words)-1], "├───", "└───", -1)
	for _, i := range words {
		fmt.Println(i)
	}
	return err
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
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
