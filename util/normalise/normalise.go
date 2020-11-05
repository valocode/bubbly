/*
Heavily inspired by (and a simplification of) <https://pkg.go.dev/k8s.io/kubernetes/pkg/kubectl/util/templates> as a strategy for normalising all printing within Bubbly to follow the same conventions.

Originally I just imported the above package, but managing kubectl/client-go/api-machinery dependencies was a total mess.

// TODO: I do not know enough about Apache License, Version 2.0 to know what is sufficient to credit. Investigate.
*/

package normalise

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
)

const Indentation = `  `

// LongDesc normalizes a command's long description to follow bubbly conventions
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normaliser{s}.heredoc().trim().string
}

// Examples normalizes a command's examples to follow bubbly conventions.
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normaliser{s}.trim().indent().string
}

type normaliser struct {
	string
}

func (s normaliser) heredoc() normaliser {
	s.string = heredoc.Doc(s.string)
	return s
}

func (s normaliser) trim() normaliser {
	s.string = strings.TrimSpace(s.string)
	return s
}

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
