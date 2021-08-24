package entfilter

import (
	"text/template"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	Extensions struct {
		entc.DefaultExtension
	}
)

func NewExtension() *Extensions {
	return &Extensions{}
}

// Templates of the extension.
func (e *Extensions) Templates() []*gen.Template {
	var templates []*gen.Template
	templates = append(
		templates,
		gen.MustParse(gen.NewTemplate("filter").
			Funcs(template.FuncMap{
				"filterNodes": filterNodes,
			}).ParseDir(
			"extensions/entfilter/templates",
		)))

	return templates
}

// filterNodes - copied from ent/contrib/entgql...
func filterNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var filteredNodes []*gen.Type
	for _, n := range nodes {
		ant := &entgql.Annotation{}
		if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes, nil
}
