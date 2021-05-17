package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
func removeFiles(files []os.FileInfo) []os.FileInfo {
	var dirs []os.FileInfo
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f)
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		if strings.Compare(dirs[i].Name(), dirs[j].Name()) > 0 {
			return false
		} else {
			return true
		}
	})
	return dirs
}

func getStrSize(q os.FileInfo) string {
	if !q.IsDir() {
		if sz := q.Size(); sz != 0 {
			return fmt.Sprintf(" (%db)", sz)
		} else {
			return " (empty)"
		}
	} else {
		return ""
	}
}

func updatePrefixes(prefix string, endLine bool) (string, string) {
	currPrefix := prefix
	newPrefix := prefix
	if endLine {
		currPrefix += "└───"
		newPrefix += "\t"
	} else {
		currPrefix += "├───"
		newPrefix += "│\t"
	}
	return currPrefix, newPrefix
}

func dirTreeWithPrefix(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := ioutil.ReadDir(path)
	if !printFiles {
		files = removeFiles(files)
	}
	for i, q := range files {
		currPrefix, newPrefix := updatePrefixes(prefix, i == len(files)-1)
		nextDir := path + string(os.PathSeparator) + q.Name()
		_, err = fmt.Fprint(out, fmt.Sprintf("%s%s%s\n", currPrefix, q.Name(), getStrSize(q)))
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
