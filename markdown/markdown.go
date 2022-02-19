package markdown

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Doc struct {
	builder *strings.Builder
}

type Table struct {
	body [][]string
}

// ListNode is a node of the list tree.
type ListNode struct {
	// Value is the content of the list item.
	Value string

	// Children is the child nodes(nested list) of the list item.
	Children []*ListNode
}

// NewDoc creates a new Markdown document struct.
func NewDoc() *Doc {
	md := new(Doc)
	md.builder = new(strings.Builder)
	return md
}

// WriteLevel1Title writes an H1 title for the given text.
func (md *Doc) WriteLevel1Title(content string) *Doc {
	md.WriteTitle(content, 1)
	return md
}

func (md *Doc) write(content string) {
	md.builder.WriteString(content)
}

// GetTitle returns header for a string with provided level.
func (md *Doc) GetTitle(content string, level int) string {
	return strings.Repeat("#", level) + " " + content
}

// WriteTitle writes header for a string with provided level.
func (md *Doc) WriteTitle(content string, level int) *Doc {
	md.write(md.GetTitle(content, level))
	md.Writeln()
	return md
}

// WriteWordLine writes a line with provided text with a newline.
func (md *Doc) WriteWordLine(content string) *Doc {
	md.Write(content)
	md.Writeln()
	return md
}

// Write writes a string to the document.
func (md *Doc) Write(content string) *Doc {
	md.write(content)
	return md
}

// Writeln writes a new line.
func (md *Doc) Writeln() *Doc {
	md.write("\n")
	return md
}

// WriteLines writes a given number of new lines.
func (md *Doc) WriteLines(lines int) *Doc {
	for i := 0; i < lines; i++ {
		md.Writeln()
	}
	return md
}

// GetMultiCode returns a multi-line code block for the given text with the given language.
func (md *Doc) GetMultiCode(content string, contentType string) string {
	return fmt.Sprintf("``` %s\n%s\n```\n", contentType, content)
}

// WriteMultiCode writes a multi-line code block for the given text with the given language.
func (md *Doc) WriteMultiCode(content, t string) *Doc {
	md.write(md.GetMultiCode(content, t))
	return md
}

// WriteCodeLine writes a single line of highlighted code for the given text..
func (md *Doc) WriteCodeLine(content string) *Doc {
	md.WriteCode(content)
	md.Writeln()
	return md
}

// GetCode returns a single line of highlighted code for the given text.
func (md *Doc) GetCode(content string) string {
	return fmt.Sprintf("`%s`", content)
}

// WriteCode writes a single line of highlighted code for the given text.
func (md *Doc) WriteCode(content string) *Doc {
	md.write(md.GetCode(content))
	return md
}

// GetLink returns a link for the given text and url.
func (md *Doc) GetLink(desc, url string) string {
	return fmt.Sprintf("[%s](%s)", desc, url)
}

// WriteLink writes a link for the given text and url.
func (md *Doc) WriteLink(desc, url string) *Doc {
	md.write(md.GetLink(desc, url))
	return md
}

// WriteLinkLine writes a link for the given text and url with a newline.
func (md *Doc) WriteLinkLine(desc, url string) *Doc {
	md.WriteLink(desc, url)
	md.WriteLines(2)
	return md
}

// GetTable returns the given table as a string.
func (md *Doc) GetTable(t *Table) string {
	return t.GetString()
}

// WriteTable writes the given table.
func (md *Doc) WriteTable(t *Table) *Doc {
	md.write(md.GetTable(t))
	return md
}

// SetTableTitle sets the title of the table on the given column.
func (t *Table) SetTableTitle(col int, content string) *Table {
	t.body[0][col] = content
	return t
}

// SetTableContent sets the content of the table on the given row and column.
func (t *Table) SetTableContent(row, col int, content string) *Table {
	row = row + 2
	t.body[row][col] = content
	return t
}

func (t *Table) GetString() string {
	var buffer strings.Builder
	for _, row := range t.body {
		buffer.WriteString("|")
		for _, col := range row {
			buffer.WriteString(col)
			buffer.WriteString("|")
		}
		buffer.WriteString("\n")

	}
	return buffer.String()
}

// NewTable constructs a blank table with the given number of rows and columns.
func NewTable(row, col int) *Table {
	t := new(Table)
	row = row + 2
	t.body = make([][]string, row)
	for i := 0; i < row; i++ {
		t.body[i] = make([]string, col)
		if i == 1 {
			for j := 0; j < col; j++ {
				t.body[i][j] = "----"
			}
		}
	}
	return t
}

// GetList return a list string from the given struct.
func GetList(tree *ListNode, level int) string {
	str := strings.Builder{}
	if tree != nil {
		str.WriteString(strings.Repeat("    ", level))
		str.WriteString("* " + tree.Value)
		str.WriteString("\n")

		for _, child := range tree.Children {
			str.WriteString(GetList(child, level+1))
		}
	} else {
		str.WriteString("\n")
	}
	return str.String()
}

// WriteList writes the given list to the document.
func (md *Doc) WriteList(tree *ListNode) *Doc {
	doc := NewDoc()
	doc.Write(GetList(tree, 0))
	return doc
}

// Export writes the entire content to a given file.
func (md *Doc) Export(filename string) error {
	return ioutil.WriteFile(filename, []byte(md.builder.String()), os.ModePerm)
}
