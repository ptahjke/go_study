package main

import (
	"bytes"
	"dir_tree/dirtree"
	"fmt"
	"os"
)

func main() {

	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirtree.DirTree(out, path, dirtree.BasicTreeDrawer{}, printFiles)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("out.String(): %v\n", out.String())
}
