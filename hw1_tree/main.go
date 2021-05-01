package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

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

func dirTreeWithPrefix(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := os.ReadDir(path)
	if !printFiles {
		var dirs []os.DirEntry
		for _, f := range files {
			if f.IsDir() {
				dirs = append(dirs, f)
			}
		}
		sort.Slice(dirs, func(i, j int) bool {
			if strings.Compare(dirs[i].Name(),dirs[j].Name()) > 0 {
				return false
			} else {
				return true
			}
		})
		files = dirs
	}
	for i, q := range files {
		currPrefix := prefix
		newPrefix := prefix
		if i == len(files)-1 {
			currPrefix += "└───"
			newPrefix += "\t"
		} else {
			currPrefix += "├───"
			newPrefix += "│\t"
		}
		nextDir := path + string(os.PathSeparator) + q.Name()
		size := ""
		if !q.IsDir() {
			if info, _ := q.Info(); info.Size() != 0 {
				size = fmt.Sprintf(" (%db)", info.Size())
			} else {
				size = " (empty)"
			}
		}
		_, err = fmt.Fprint(out, fmt.Sprintf("%s%s%s\n", currPrefix, q.Name(), size))
		if err != nil {
			return err
		}
		if q.IsDir() {
			err = dirTreeWithPrefix(out, nextDir, printFiles, newPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeWithPrefix(out, path, printFiles, "")
}
