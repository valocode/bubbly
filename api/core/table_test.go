package core

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

var productTable = Table{
	Name: "product",
	Fields: []TableField{
		{
			Name:   "id",
			Type:   cty.String,
			Unique: true,
		},
	},
}

var projectTable = Table{
	Name: "project",
	Fields: []TableField{
		{
			Name:   "id",
			Type:   cty.String,
			Unique: true,
		},
		{
			Name:   "product_id",
			Type:   cty.List(cty.String),
			Unique: false,
		},
	},
}

var repoTable = Table{
	Name: "repository",
	Fields: []TableField{
		{
			Name:   "id",
			Type:   cty.String,
			Unique: true,
		},
		{
			Name:   "url",
			Type:   cty.String,
			Unique: false,
		},
		{
			Name:   "project_id",
			Type:   cty.List(cty.String),
			Unique: false,
		},
	},
}

var repoVersionTable = Table{
	Name: "repository_version",
	Fields: []TableField{
		{
			Name:   "id",
			Type:   cty.String,
			Unique: true,
		},
		{
			Name:   "tag",
			Type:   cty.String,
			Unique: false,
		},
		{
			Name:   "branch",
			Type:   cty.String,
			Unique: false,
		},
		{
			Name:   "repository_id",
			Type:   cty.String,
			Unique: false,
		},
	},
}

func TestIssueTable(t *testing.T) {
	issueTable := Table{
		Name: "linter_issue",
		Fields: []TableField{
			{
				Name:   "id",
				Type:   cty.String,
				Unique: true,
			},
			{
				Name:   "repository_version_id",
				Type:   cty.String,
				Unique: true,
			},
			{
				Name:   "severity",
				Type:   cty.String,
				Unique: true,
			},
			{
				Name:   "type",
				Type:   cty.String,
				Unique: true,
			},
		},
	}

	t.Logf("Is there a test needed here for table %s?", issueTable.Name)
}
