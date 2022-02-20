package markdown

import "strings"

type Table struct {
	body [][]string
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

// GetTable returns the given table as a string.
func (t *Table) GetTable() string {
	return t.GetTableString()
}

// GetTableString returns the given table as a string.
func (t *Table) GetTableString() string {
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
