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

// WriteLevel1Title writes an H1 title for the given text.
func (md *Doc) WriteLevel1Title(content string) (*Doc, error) {
	_, err := md.WriteHeader(content, 1)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func (md *Doc) write(content string) error {
	_, err := md.builder.WriteString(content)
	return err
}

// WriteHeader writes header for a string with provided level.
func (md *Doc) WriteHeader(content string, level int) (*Doc, error) {
	err := md.write(GetHeader(content, level))
	if err != nil {
		return nil, err
	}
	_, err = md.Writeln()
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteWordLine writes a line with provided text with a newline.
func (md *Doc) WriteWordLine(content string) (*Doc, error) {
	err := md.write(content)
	if err != nil {
		return nil, err
	}
	_, err = md.Writeln()
	if err != nil {
		return nil, err
	}
	return md, nil
}

// Write writes a string to the document.
func (md *Doc) Write(content string) (*Doc, error) {
	err := md.write(content)
	if err != nil {
		return nil, err
	}
	return md, nil
}

// Writeln writes a new line.
func (md *Doc) Writeln() (*Doc, error) {
	err := md.write("\n")
	if err != nil {
		return nil, err
	}
	return md, err
}

// WriteLines writes a given number of new lines.
func (md *Doc) WriteLines(lines int) (*Doc, error) {
	for i := 0; i < lines; i++ {
		_, err := md.Writeln()
		if err != nil {
			return nil, err
		}
	}
	return md, nil
}

// WriteMultiCode writes a multi-line code block for the given text with the given language.
func (md *Doc) WriteMultiCode(content, t string) (*Doc, error) {
	err := md.write(GetMultiCode(content, t))
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteCodeLine writes a single line of highlighted code for the given text..
func (md *Doc) WriteCodeLine(content string) (*Doc, error) {
	_, err := md.WriteCode(content)
	if err != nil {
		return nil, err
	}
	_, err = md.Writeln()
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteCode writes a single line of highlighted code for the given text.
func (md *Doc) WriteCode(content string) (*Doc, error) {
	err := md.write(GetMonospaceCode(content))
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteLink writes a link for the given text and url.
func (md *Doc) WriteLink(desc, url string) (*Doc, error) {
	err := md.write(GetLink(desc, url))
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteLinkLine writes a link for the given text and url with a newline.
func (md *Doc) WriteLinkLine(desc, url string) (*Doc, error) {
	_, err := md.WriteLink(desc, url)
	if err != nil {
		return nil, err
	}
	_, err = md.WriteLines(2)
	if err != nil {
		return nil, err
	}
	return md, nil
}

// WriteTable writes the given table.
func (md *Doc) WriteTable(t *Table) (*Doc, error) {
	err := md.write(t.GetTable())
	if err != nil {
		return nil, err
	}
	return md, err
}

// WriteList writes the given list to the document.
func (md *Doc) WriteList(tree *ListNode) (*Doc, error) {
	_, err := md.Write(tree.GetList(0))
	if err != nil {
		return nil, err
	}
	return md, err
}

// Export writes the entire content to a given file.
func (md *Doc) Export(filename string) error {
	return ioutil.WriteFile(filename, []byte(md.builder.String()), os.ModePerm)
}
