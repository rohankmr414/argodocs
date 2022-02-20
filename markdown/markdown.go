package markdown

import (
	"io/ioutil"
	"os"
	"strings"
)

type Doc struct {
	builder *strings.Builder
}

// NewDoc creates a new Markdown document struct.
func NewDoc() *Doc {
	md := new(Doc)
	md.builder = new(strings.Builder)
	return md
}

// write appends the given string to the document.
func (md *Doc) write(content string) {
	md.builder.WriteString(content)
}

// WriteHeader writes header for a string with provided level.
func (md *Doc) WriteHeader(content string, level int) *Doc {
	md.write(GetHeader(content, level))
	md.Writeln()

	return md
}

// Write writes a string to the document.
func (md *Doc) Write(content string) *Doc {
	md.write(content)

	return md
}

// Writeln writes a string to the document and appends a newline.
func (md *Doc) Writeln(content string) *Doc {
	md.write(content + "\n")
	return md
}

// WriteLines writes a given number of new lines.
func (md *Doc) WriteLines(lines int) *Doc {
	for i := 0; i < lines; i++ {
		md.write("\n")
	}
	return md
}

// WriteMultiCode writes a multi-line code block for the given text with the given language.
func (md *Doc) WriteMultiCode(content, t string) *Doc {
	md.write(GetMultiCode(content, t))
	return md
}

// WriteCode writes a single line of highlighted code for the given text.
func (md *Doc) WriteCode(content string) *Doc {
	md.write(GetMonospaceCode(content))
	return md
}

// WriteLink writes a link for the given text and url.
func (md *Doc) WriteLink(desc, url string) *Doc {
	md.write(GetLink(desc, url))
	return md
}

// WriteTable writes the given table.
func (md *Doc) WriteTable(t *Table) *Doc {
	md.write(t.GetTable())
	return md
}

// WriteList writes the given list to the document.
func (md *Doc) WriteList(tree *ListNode) *Doc {
	md.Write(tree.GetList(0, -1))
	return md
}

// Export writes the entire content to a given file.
func (md *Doc) Export(filename string) error {
	return ioutil.WriteFile(filename, []byte(md.builder.String()), os.ModePerm)
}
