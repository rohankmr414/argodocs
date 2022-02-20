package markdown

import (
	"fmt"
	"strings"
)

// ListType is an enumeration of list types.
type ListType int

const (
	// ListTypeOrdered is an ordered list.
	ListTypeOrdered ListType = iota

	// ListTypeUnordered is an unordered list.
	ListTypeUnordered
)

// ListNode is a node of the list tree.
type ListNode struct {
	// Value is the content of the list item.
	Value string

	// NodeType is the type of the parent list(ListTypeOrdered | ListTypeUnordered).
	NodeType ListType

	// Children is the child nodes(nested list) of the list item.
	Children []*ListNode
}

// GetList return a list string from the given struct.
func (tree *ListNode) GetList(count, level int) string {
	str := strings.Builder{}
	if tree != nil {
		if tree.Value != "" {
			str.WriteString(strings.Repeat("    ", level))
			if tree.NodeType == ListTypeOrdered {
				li := fmt.Sprintf("%d. %s\n", count, tree.Value)
				str.WriteString(li)
			} else if tree.NodeType == ListTypeUnordered {
				li := fmt.Sprintf("- %s\n", tree.Value)
				str.WriteString(li)
			}
		}

		for i, child := range tree.Children {
			str.WriteString(child.GetList(i+1, level+1))
		}
	}

	return str.String()
}
