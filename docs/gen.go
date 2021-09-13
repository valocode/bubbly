package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/spf13/cobra/doc"
	"github.com/valocode/bubbly/cmd"
	"github.com/valocode/bubbly/env"
)

func main() {
	genCLIDocs()
	genSchemaDocs()
}

func genCLIDocs() {
	cmd := cmd.NewCmdRoot(env.NewBubblyConfig())
	err := doc.GenMarkdownTreeCustom(cmd, "./docs/docs/cli",
		func(s string) string {
			filename := filepath.Base(s)
			name := filename[:len(filename)-len(filepath.Ext(filename))]
			name = strings.Join(strings.Split(name, "_"), " ")
			return fmt.Sprintf(`---
title: %s
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- cli
---
`, name)
		},
		func(s string) string { return s },
	)
	// err := doc.GenMarkdownTree(cmd, "./docs/docs/cli")
	if err != nil {
		log.Fatal(err)
	}
}

func genSchemaDocs() {
	graph, err := entc.LoadGraph("./ent/schema", &gen.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, graph); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./docs/docs/schema/schema.md", b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}

var tmpl = template.Must(template.New("schema").
	Funcs(template.FuncMap{
		"fmtType": func(s string) string {
			return strings.NewReplacer(
				".", "DOT",
				"*", "STAR",
				"[", "LBRACK",
				"]", "RBRACK",
			).Replace(s)
		},
		"toLower": func(s string) string {
			return strings.ToLower(s)
		},
	}).
	Parse(`---
title: Schema
hide_title: false
hide_table_of_contents: false
keywords:
- docs
- bubbly
- schema
---

import Mermaid from '@theme/Mermaid';

## Introduction

The schema is an core part of bubbly - it is needed to store all the relevant Release Readiness data over which policies are run.

The library used to define the schema is [entgo](https://entgo.io/) which is an amazing framework, and this documentation was generated using it :)

## Overview

The following diagram shows the bubbly schema
{{- with $.Nodes }}
<Mermaid chart={` + "`" + `
erDiagram
{{- range $n := . }}
    {{ $n.Name }} {
        {{ fmtType $n.ID.Type.String }} {{ $n.ID.Name }}
	{{- range $f := $n.Fields }}
        {{ fmtType $f.Type.String }} {{ $f.Name }}
	{{- end }}
    }
{{- end }}
{{- range $n := . }}
    {{- range $e := $n.Edges }}
	{{- if not $e.IsInverse }}
		{{- $rt := "|o--o|" }}{{ if $e.O2M }}{{ $rt = "|o--o{" }}{{ else if $e.M2O }}{{ $rt = "}o--o|" }}{{ else if $e.M2M }}{{ $rt = "}o--o{" }}{{ end }}
    	{{ $n.Name }} {{ $rt }} {{ $e.Type.Name }} : "{{ $e.Name }}{{ with $e.Ref }}/{{ .Name }}{{ end }}"
	{{- end }}
	{{- end }}
{{- end }}
` + "`" + `}/>

## Types

The types in the schema (AKA tables in the SQL schema) are listed below.
This list is auto-generated from the ent schema.

{{- range $n := . }}
### {{ $n.Name }}

#### Fields
{{- range $f := $n.Fields }}
- **{{ $f.Name }}** ({{ $f.Type }})
{{- end }}

#### Edges
{{- range $e := $n.Edges }}
- **{{ $e.Name }}** ({{ $e.Rel.Type }} to [{{ $e.Type.Name }}](#{{ toLower $e.Type.Name }}))
{{- end }}

{{- end }}
{{- end }}
`))
