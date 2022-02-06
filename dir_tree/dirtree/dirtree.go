package dirtree

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

type TreeDrawer interface {
	Draw(*bytes.Buffer, Node)
}

type Node struct {
	name   string
	nodes  []Node
	isLast bool
	size   int64
	isFile bool
}

func DirTree(out *bytes.Buffer, path string, drawer TreeDrawer, withFiles bool) (err error) {
	rootNode := Node{}
	err = recursiveFillDirectories(path, &rootNode, withFiles)
	if err != nil {
		return err
	}

	drawer.Draw(out, rootNode)

	return nil
}

func recursiveFillDirectories(path string, parentNode *Node, withFiles bool) error {
	dirs, err := getDirs(path, withFiles)
	if err != nil {
		return err
	}

	for i, dir := range dirs {
		nextPath := path + string(os.PathSeparator) + dir.Name()
		finfo, err := os.Stat(nextPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var node Node = Node{
			name:   dir.Name(),
			isLast: i == len(dirs)-1,
		}

		if finfo.IsDir() {
			recursiveFillDirectories(nextPath, &node, withFiles)
			parentNode.nodes = append(parentNode.nodes, node)
		} else if withFiles {
			node.size = finfo.Size()
			node.isFile = true
			parentNode.nodes = append(parentNode.nodes, node)
		}
	}

	return nil
}

func getDirs(path string, withFiles bool) ([]os.FileInfo, error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if !withFiles {
		var res []os.FileInfo
		for _, dir := range dirs {
			if dir.IsDir() {
				res = append(res, dir)
			}
		}

		return res, nil
	}

	return dirs, nil
}
