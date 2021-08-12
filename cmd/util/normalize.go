package util

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
)

const Indentation = `  `

type normaliser struct {
	string
}

// LongDesc normalizes a command's long description to follow bubbly conventions
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normaliser{s}.heredoc().trim().string
}

// Examples normalizes a command's examples to follow bubbly conventions. I.e.,
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normaliser{s}.trim().indent().string
}

// returns a here document from given string
func (s normaliser) heredoc() normaliser {
	s.string = heredoc.Doc(s.string)
	return s
}

// removes trailing/leading whitespaces from a string
func (s normaliser) trim() normaliser {
	s.string = strings.TrimSpace(s.string)
	return s
}

// indents a string according to the pre-defined Indentation
func (s normaliser) indent() normaliser {
	indentedLines := []string{}
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}
