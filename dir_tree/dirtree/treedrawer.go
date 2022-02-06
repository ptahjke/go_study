package dirtree

import (
	"bytes"
	"fmt"
)

type BasicTreeDrawer struct {
}

func (drawer BasicTreeDrawer) Draw(out *bytes.Buffer, node Node) {
	drawDirectoriesTree(out, node)
}

func drawDirectoriesTree(out *bytes.Buffer, node Node) {
	res := recursiveDrawDirectoriesTree(node, 0, 0)
	out.WriteString(res)
}

const commonElementSeparator = "├───"
const lastElementSeparator = "└───"
const vericalDerictoriesElementSeparator = "│"
const horizontalElementSeparator = "	"

func recursiveDrawDirectoriesTree(node Node, offset int, lastParentCount int) (res string) {
	for _, node := range node.nodes {
		toffset := offset
		tlastParentCount := lastParentCount
		for ; toffset > 0; toffset-- {
			if tlastParentCount != toffset {
				res += vericalDerictoriesElementSeparator
			} else {
				tlastParentCount--
			}

			res += horizontalElementSeparator
		}

		if node.isLast {
			lastParentCount++
			res += lastElementSeparator
		} else {
			res += commonElementSeparator
		}

		res += node.name

		if node.isFile {
			size := "empty"
			if node.size > 0 {
				size = fmt.Sprint(node.size) + "b"
			}
			res += " (" + size + ")"
		}

		res += "\n"

		if len(node.nodes) > 0 {
			rres := recursiveDrawDirectoriesTree(node, offset+1, lastParentCount)

			res += rres
		}
	}

	return res
}
