package markdown

import (
	"fmt"
	"strings"
)

type ListType int

const (
	ListTypeOrdered ListType = iota
	ListTypeUnordered
)

// ListNode is a node of the list tree.
type ListNode struct {
	// Value is the content of the list item.
	Value string

	// ChildrenType is the type of the Children list(ListTypeOrdered | ListTypeUnordered).
	ChildrenType ListType

	// Children is the child nodes(nested list) of the list item.
	Children []*ListNode
}

// GetList return a list string from the given struct.
func (tree *ListNode) GetList(count, level int) string {
	str := strings.Builder{}
	if tree != nil {
		if tree.Value == "" {
			str.WriteString("\n")
		} else {
			str.WriteString(strings.Repeat("    ", level))
			if tree.ChildrenType == ListTypeOrdered {
				li := fmt.Sprintf("%d. %s\n", count, tree.Value)
				str.WriteString(li)
			} else {
				li := fmt.Sprintf("* %s\n", tree.Value)
				str.WriteString(li)
			}
		}

		for i, child := range tree.Children {
			str.WriteString(child.GetList(i+1, level+1))
		}

	} else {
		str.WriteString("\n")
	}
	return str.String()
}
