package builtin

import (
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
)

var BuiltinTables = []core.Table{

	// #######################################
	// _SCHEMA
	// #######################################
	table("_schema",
		fields(
			field("tables", cty.Map(cty.DynamicPseudoType), false),
		),
		joins(),
	),
	// #######################################
	// _RESOURCE
	// #######################################
	table("_resource",
		fields(
			field("id", cty.String, true),
			field("name", cty.String, false),
			field("kind", cty.String, false),
			field("api_version", cty.String, false),
			field("spec", cty.String, false),
			field("metadata", cty.Object(map[string]cty.Type{"labels": cty.Map(cty.String)}), false),
		),
		joins(),
	),
	// #######################################
	// _EVENT
	// #######################################
	table("_event",
		fields(
			field("status", cty.String, false),
			field("error", cty.String, false),
			field("time", parser.TimeType, false),
		),
		joins(
			join("_resource", false, false),
		),
	),
	// #######################################
	// RELEASE_ENTRY
	// #######################################
	table("release_entry",
		fields(
			field("name", cty.String, false),
			field("result", cty.Bool, false),
			field("reason", cty.String, false),
		),
		joins(
			join("release", false, false),
			join("release_criteria", false, false),
			join("_resource", false, false),
		),
	),
	// #######################################
	// RELEASE
	// #######################################
	table("release",
		fields(
			field("name", cty.String, true),
			field("version", cty.String, true),
		),
		joins(
			join("project", false, true),
		),
	),
	// #######################################
	// RELEASE_STAGE
	// #######################################
	table("release_stage",
		fields(
			field("name", cty.String, true),
		),
		joins(
			join("release", false, true),
		),
	),
	// #######################################
	// RELEASE_CRITERIA
	// #######################################
	table("release_criteria",
		fields(
			field("entry_name", cty.String, true),
		),
		joins(
			join("release_stage", false, false),
			join("release", false, true),
		),
	),
	// #######################################
	// CODE_SCAN
	// #######################################
	table("code_scan",
		fields(
			field("tool", cty.String, false),
		),
		joins(
			join("release", false, false),
			join("lifecycle", true, false),
		),
	),
	// #######################################
	// LIFECYCLE
	// #######################################
	table("lifecycle",
		fields(
			field("status", cty.String, false),
		),
		joins(),
	),
	// #######################################
	// LIFECYCLE_ENTRY
	// #######################################
	table("lifecycle_entry",
		fields(
			field("author", cty.String, false),
			field("message", cty.String, false),
			field("signed", cty.Bool, false),
		),
		joins(
			join("lifecycle", false, false),
		),
	),
	// #######################################
	// TEST_RUN
	// #######################################
	table("test_run",
		fields(
			field("tool", cty.String, false),
			field("type", cty.String, false),
			field("name", cty.String, false),
			field("elapsed", cty.Number, false),
			field("result", cty.Bool, false),
		),
		joins(
			join("release", false, false),
			join("lifecycle", true, false),
		),
	),
	// #######################################
	// TEST_CASE
	// #######################################
	table("test_case",
		fields(
			field("name", cty.String, false),
			field("result", cty.Bool, false),
			field("message", cty.String, false),
		),
		joins(
			join("test_run", false, false),
		),
	),
	// #######################################
	// CODE_ISSUE
	// #######################################
	table("code_issue",
		fields(
			field("id", cty.String, false),
			field("message", cty.String, false),
			field("severity", cty.String, false),
			field("type", cty.String, false),
		),
		joins(
			join("code_scan", false, false),
		),
	),
	// #######################################
	// PROJECT
	// #######################################
	table("project",
		fields(
			field("name", cty.String, true),
		),
		joins(),
	),
	// #######################################
	// REPO
	// #######################################
	table("repo",
		fields(
			field("id", cty.String, true),
			field("name", cty.String, false),
		),
		joins(
			join("project", false, false),
		),
	),
	// #######################################
	// BRANCH
	// #######################################
	table("branch",
		fields(
			field("name", cty.String, true),
		),
		joins(
			join("repo", false, true),
		),
	),
	// #######################################
	// COMMIT
	// #######################################
	table("commit",
		fields(
			field("id", cty.String, true),
			field("tag", cty.String, false),
			field("time", cty.String, false),
		),
		joins(
			join("branch", false, false),
			join("repo", false, true),
		),
	),
	// #######################################
	// RELEASE_ITEM
	// #######################################
	table("release_item",
		fields(
			field("type", cty.String, false),
		),
		joins(
			join("release", false, false),
			join("commit", true, true),
			join("artifact", true, true),
		),
	),
	// #######################################
	// ARTIFACT
	// #######################################
	table("artifact",
		fields(
			field("name", cty.String, false),
			field("sha256", cty.String, true),
			field("location", cty.String, false),
		),
		joins(),
	),
}
