package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFileSize1(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

type Tuple struct {
	Name    string
	Index   int
	path    string
	symbols string
}

type Dirs struct {
	Items []Tuple
}

func (arr *Dirs) AddItem(item Tuple) []Tuple {
	arr.Items = append(arr.Items, item)
	return arr.Items
}

func (arr *Dirs) CheckValue(item Tuple) bool {
	for _, val := range arr.Items {
		if val.Name == item.Name && val.Index == item.Index && val.path == item.path {
			return true
		}
	}
	return false
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var space string
	var fileSize string
	items := []Tuple{}
	dirs := Dirs{items}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if printFiles {
			byteSize, _ := GetFileSize1(path)
			if byteSize == 0 {
				fileSize = " (empty)"
			} else {
				fileSize = " (" + fmt.Sprintf("%v", byteSize) + "b)"
			}
		}
		if !info.IsDir() {
			for idx, f := range strings.Split(path+fileSize, "/")[1:] {
				if strings.ContainsAny(f, ".") && !printFiles {
					continue
				}
				item := Tuple{Name: f, Index: idx, path: strings.Split(path+fileSize, f)[0]}
				if dirs.CheckValue(item) {
					continue
				} else {
					dirs.AddItem(item)
				}
			}
		}
		return err
	})
	for idx, item := range dirs.Items {
		space = strings.Repeat("        ", item.Index)
		for i := idx + 1; i < len(dirs.Items); i++ {
			if item.Index == dirs.Items[i].Index && item.path == dirs.Items[i].path {
				dirs.Items[idx].symbols = space + "├───" + item.Name
				break
			} else {
				dirs.Items[idx].symbols = space + "└───" + item.Name
			}
		}
	}
	for idx, item := range dirs.Items {
		char := strings.Index(item.symbols, "├")
		if char > -1 {
			for i := 0; i < len(dirs.Items[idx+1:]); i++ {
				if len(string([]rune(dirs.Items[idx+1:][i].symbols))) != 0 && string([]rune(dirs.Items[idx+1:][i].symbols)[char]) == " " {
					dirs.Items[idx+1:][i].symbols =
						string([]rune(dirs.Items[idx+1:][i].symbols)[:char]) + "|" + string([]rune(dirs.Items[idx+1:][i].symbols)[char+1:])
				} else {
					break
				}
			}
		}
	}
	if printFiles {
		for i := len(dirs.Items) - 1; i >= 0; i-- {
			if len(string([]rune(dirs.Items[i].symbols))) != 0 && string([]rune(dirs.Items[i].symbols)[0]) == "├" {
				dirs.Items[i].symbols = "└" + string([]rune(dirs.Items[i].symbols)[1:])
				break
			} else if len(string([]rune(dirs.Items[i].symbols))) != 0 && string([]rune(dirs.Items[i].symbols)[0]) == "|" {
				dirs.Items[i].symbols = " " + string([]rune(dirs.Items[i].symbols)[1:])
			}
		}
	}
	for _, item := range dirs.Items {
		fmt.Fprintln(out, item.symbols)
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
