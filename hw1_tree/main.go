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
	// get the size
	return fi.Size(), nil
}

type Tuple struct {
	Name  string
	Index int
	path  string
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
	//var dirs []string
	var count int
	var tree []string
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
		fmt.Println(path)
		if !info.IsDir() {
			for idx, f := range strings.Split(path+fileSize, "/")[1:] {
				item := Tuple{Name: f, Index: idx + 1, path: strings.Split(path+fileSize, f)[0]}
				if dirs.CheckValue(item) {
					continue
				} else {
					dirs.AddItem(item)
					//fmt.Println(item)
				}
				count += 1
			}
			//for _, item := range dirs.Items {
			//	space = strings.Repeat("│       ", item.Index)
			//	tree = append(tree, space + "├───" + item.Name + "\n")
			//}
		}

		for _, item := range dirs.Items {
			space = strings.Repeat("│       ", item.Index)
			tree = append(tree, space+"├───"+item.Name+"\n")
		}
		return err
	})
	for idx, str := range tree {
		if idx == 0 {
			continue
		} else {
			if len(strings.Split(str, "├───")[0]) < len(strings.Split(tree[idx-1], "├───")[0]) {
				tree[idx-1] = strings.Replace(tree[idx-1], "├───", "└───", -1)
			}
		}
	}

	//tree[len(tree)-1] = strings.Replace(tree[len(tree)-1], "├───", "└───", -1)
	for _, i := range tree[len(tree)-count:] {
		fmt.Fprint(out, i[10:])
	}
	return err
}

//		for idx, f := range strings.Split(path+fileSize, "/") {
//			//if idx == 0 {
//			//	continue
//			//}
//			space = strings.Repeat("│       ", idx)
//			if contains(words, space+"├───"+f) {
//				continue
//			} else {
//				words = append(words, space+"├───"+f)
//			}
//		}
//	}
//	return nil
//})
//var prevLen = 0
//var prevStr string
//for idx, str := range words {
//	result := strings.Split(str, "├───")
//	if prevLen > len(result[0]) {
//		words[idx-1] = strings.Replace(prevStr, "├───", "└───", -1)
//		prevLen = 0
//	} else {
//		prevStr = str
//		prevLen = len(result[0])
//	}
//}
//words[len(words)-1] = strings.Replace(words[len(words)-1], "├───", "└───", -1)
//var tree string
//if printFiles {
//	for _, i := range words {
//		tree = tree + "\n " + i[5:]
//	}
//} else {
//	for _, i := range words {
//		if strings.ContainsAny(i[5:], ".") {
//			continue
//		} else {
//			tree = tree + "\n " + i[5:]
//		}
//	}
//}
//fmt.Fprintln(out, tree)
////fmt.Println(tree)
//	return err
//}
//

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
