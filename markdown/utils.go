package markdown

import (
	"fmt"
	"strings"
)

// GetHeader returns header for a string with provided level.
func GetHeader(content string, level int) string {
	return strings.Repeat("#", level) + " " + content
}

// GetMultiCode returns a multi-line code block for the given text with the given language.
func GetMultiCode(contentType, content string) string {
	return fmt.Sprintf("``` %s\n%s\n```\n", contentType, content)
}

// GetMonospaceCode returns a single line of monospace highlighted code for the given text.
func GetMonospaceCode(content string) string {
	return fmt.Sprintf("`%s`", content)
}

// GetLink returns a link for the given text and url.
func GetLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}
