package workflow

import "strings"

func cleanupComment(comment string) string {
	trimmedCommentLines := removeLeadingPoundSigns(strings.Split(comment, "\n"))
	return strings.Join(trimmedCommentLines, "\n")
}

func removeLeadingPoundSigns(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmedLine := strings.TrimLeftFunc(line, func(r rune) bool {
			return r == '#' || r == ' '
		})

		result = append(result, trimmedLine)
	}

	return result
}
