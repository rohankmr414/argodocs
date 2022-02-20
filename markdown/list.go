package markdown

import "strings"

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
func (tree *ListNode) GetList(level int) string {
	str := strings.Builder{}
	if tree != nil {
		str.WriteString(strings.Repeat("    ", level))
		str.WriteString("* " + tree.Value)
		str.WriteString("\n")

		for _, child := range tree.Children {
			str.WriteString(child.GetList(level + 1))
		}
	} else {
		str.WriteString("\n")
	}
	return str.String()
}
