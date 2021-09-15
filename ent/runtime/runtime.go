// Code generated by entc, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/valocode/bubbly/ent/adapter"
	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/codescan"
	"github.com/valocode/bubbly/ent/component"
	"github.com/valocode/bubbly/ent/event"
	"github.com/valocode/bubbly/ent/gitcommit"
	"github.com/valocode/bubbly/ent/license"
	"github.com/valocode/bubbly/ent/organization"
	"github.com/valocode/bubbly/ent/project"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releasecomponent"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/ent/releasepolicy"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	"github.com/valocode/bubbly/ent/repo"
	"github.com/valocode/bubbly/ent/schema"
	"github.com/valocode/bubbly/ent/testcase"
	"github.com/valocode/bubbly/ent/testrun"
	"github.com/valocode/bubbly/ent/vulnerability"
	"github.com/valocode/bubbly/ent/vulnerabilityreview"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	adapterFields := schema.Adapter{}.Fields()
	_ = adapterFields
	// adapterDescName is the schema descriptor for name field.
	adapterDescName := adapterFields[0].Descriptor()
	// adapter.NameValidator is a validator for the "name" field. It is called by the builders before save.
	adapter.NameValidator = func() func(string) error {
		validators := adapterDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// adapterDescTag is the schema descriptor for tag field.
	adapterDescTag := adapterFields[1].Descriptor()
	// adapter.TagValidator is a validator for the "tag" field. It is called by the builders before save.
	adapter.TagValidator = adapterDescTag.Validators[0].(func(string) error)
	// adapterDescModule is the schema descriptor for module field.
	adapterDescModule := adapterFields[2].Descriptor()
	// adapter.ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	adapter.ModuleValidator = adapterDescModule.Validators[0].(func(string) error)
	artifactHooks := schema.Artifact{}.Hooks()
	artifact.Hooks[0] = artifactHooks[0]
	artifactFields := schema.Artifact{}.Fields()
	_ = artifactFields
	// artifactDescName is the schema descriptor for name field.
	artifactDescName := artifactFields[0].Descriptor()
	// artifact.NameValidator is a validator for the "name" field. It is called by the builders before save.
	artifact.NameValidator = artifactDescName.Validators[0].(func(string) error)
	// artifactDescSha256 is the schema descriptor for sha256 field.
	artifactDescSha256 := artifactFields[1].Descriptor()
	// artifact.Sha256Validator is a validator for the "sha256" field. It is called by the builders before save.
	artifact.Sha256Validator = artifactDescSha256.Validators[0].(func(string) error)
	// artifactDescTime is the schema descriptor for time field.
	artifactDescTime := artifactFields[3].Descriptor()
	// artifact.DefaultTime holds the default value on creation for the time field.
	artifact.DefaultTime = artifactDescTime.Default.(func() time.Time)
	codeissueFields := schema.CodeIssue{}.Fields()
	_ = codeissueFields
	// codeissueDescRuleID is the schema descriptor for rule_id field.
	codeissueDescRuleID := codeissueFields[0].Descriptor()
	// codeissue.RuleIDValidator is a validator for the "rule_id" field. It is called by the builders before save.
	codeissue.RuleIDValidator = codeissueDescRuleID.Validators[0].(func(string) error)
	// codeissueDescMessage is the schema descriptor for message field.
	codeissueDescMessage := codeissueFields[1].Descriptor()
	// codeissue.MessageValidator is a validator for the "message" field. It is called by the builders before save.
	codeissue.MessageValidator = codeissueDescMessage.Validators[0].(func(string) error)
	codescanHooks := schema.CodeScan{}.Hooks()
	codescan.Hooks[0] = codescanHooks[0]
	codescanFields := schema.CodeScan{}.Fields()
	_ = codescanFields
	// codescanDescTool is the schema descriptor for tool field.
	codescanDescTool := codescanFields[0].Descriptor()
	// codescan.ToolValidator is a validator for the "tool" field. It is called by the builders before save.
	codescan.ToolValidator = codescanDescTool.Validators[0].(func(string) error)
	// codescanDescTime is the schema descriptor for time field.
	codescanDescTime := codescanFields[1].Descriptor()
	// codescan.DefaultTime holds the default value on creation for the time field.
	codescan.DefaultTime = codescanDescTime.Default.(func() time.Time)
	componentFields := schema.Component{}.Fields()
	_ = componentFields
	// componentDescName is the schema descriptor for name field.
	componentDescName := componentFields[0].Descriptor()
	// component.NameValidator is a validator for the "name" field. It is called by the builders before save.
	component.NameValidator = componentDescName.Validators[0].(func(string) error)
	// componentDescVendor is the schema descriptor for vendor field.
	componentDescVendor := componentFields[1].Descriptor()
	// component.DefaultVendor holds the default value on creation for the vendor field.
	component.DefaultVendor = componentDescVendor.Default.(string)
	// componentDescVersion is the schema descriptor for version field.
	componentDescVersion := componentFields[2].Descriptor()
	// component.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	component.VersionValidator = componentDescVersion.Validators[0].(func(string) error)
	eventFields := schema.Event{}.Fields()
	_ = eventFields
	// eventDescMessage is the schema descriptor for message field.
	eventDescMessage := eventFields[0].Descriptor()
	// event.DefaultMessage holds the default value on creation for the message field.
	event.DefaultMessage = eventDescMessage.Default.(string)
	// eventDescTime is the schema descriptor for time field.
	eventDescTime := eventFields[2].Descriptor()
	// event.DefaultTime holds the default value on creation for the time field.
	event.DefaultTime = eventDescTime.Default.(func() time.Time)
	gitcommitFields := schema.GitCommit{}.Fields()
	_ = gitcommitFields
	// gitcommitDescHash is the schema descriptor for hash field.
	gitcommitDescHash := gitcommitFields[0].Descriptor()
	// gitcommit.HashValidator is a validator for the "hash" field. It is called by the builders before save.
	gitcommit.HashValidator = gitcommitDescHash.Validators[0].(func(string) error)
	// gitcommitDescBranch is the schema descriptor for branch field.
	gitcommitDescBranch := gitcommitFields[1].Descriptor()
	// gitcommit.BranchValidator is a validator for the "branch" field. It is called by the builders before save.
	gitcommit.BranchValidator = gitcommitDescBranch.Validators[0].(func(string) error)
	licenseFields := schema.License{}.Fields()
	_ = licenseFields
	// licenseDescSpdxID is the schema descriptor for spdx_id field.
	licenseDescSpdxID := licenseFields[0].Descriptor()
	// license.SpdxIDValidator is a validator for the "spdx_id" field. It is called by the builders before save.
	license.SpdxIDValidator = licenseDescSpdxID.Validators[0].(func(string) error)
	// licenseDescName is the schema descriptor for name field.
	licenseDescName := licenseFields[1].Descriptor()
	// license.NameValidator is a validator for the "name" field. It is called by the builders before save.
	license.NameValidator = licenseDescName.Validators[0].(func(string) error)
	// licenseDescIsOsiApproved is the schema descriptor for is_osi_approved field.
	licenseDescIsOsiApproved := licenseFields[4].Descriptor()
	// license.DefaultIsOsiApproved holds the default value on creation for the is_osi_approved field.
	license.DefaultIsOsiApproved = licenseDescIsOsiApproved.Default.(bool)
	organizationFields := schema.Organization{}.Fields()
	_ = organizationFields
	// organizationDescName is the schema descriptor for name field.
	organizationDescName := organizationFields[0].Descriptor()
	// organization.NameValidator is a validator for the "name" field. It is called by the builders before save.
	organization.NameValidator = organizationDescName.Validators[0].(func(string) error)
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectFields[0].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	releaseHooks := schema.Release{}.Hooks()
	release.Hooks[0] = releaseHooks[0]
	releaseFields := schema.Release{}.Fields()
	_ = releaseFields
	// releaseDescName is the schema descriptor for name field.
	releaseDescName := releaseFields[0].Descriptor()
	// release.NameValidator is a validator for the "name" field. It is called by the builders before save.
	release.NameValidator = releaseDescName.Validators[0].(func(string) error)
	// releaseDescVersion is the schema descriptor for version field.
	releaseDescVersion := releaseFields[1].Descriptor()
	// release.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	release.VersionValidator = releaseDescVersion.Validators[0].(func(string) error)
	releasecomponentHooks := schema.ReleaseComponent{}.Hooks()
	releasecomponent.Hooks[0] = releasecomponentHooks[0]
	releasecomponentFields := schema.ReleaseComponent{}.Fields()
	_ = releasecomponentFields
	releaseentryFields := schema.ReleaseEntry{}.Fields()
	_ = releaseentryFields
	// releaseentryDescTime is the schema descriptor for time field.
	releaseentryDescTime := releaseentryFields[1].Descriptor()
	// releaseentry.DefaultTime holds the default value on creation for the time field.
	releaseentry.DefaultTime = releaseentryDescTime.Default.(func() time.Time)
	releasepolicyFields := schema.ReleasePolicy{}.Fields()
	_ = releasepolicyFields
	// releasepolicyDescName is the schema descriptor for name field.
	releasepolicyDescName := releasepolicyFields[0].Descriptor()
	// releasepolicy.NameValidator is a validator for the "name" field. It is called by the builders before save.
	releasepolicy.NameValidator = func() func(string) error {
		validators := releasepolicyDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// releasepolicyDescModule is the schema descriptor for module field.
	releasepolicyDescModule := releasepolicyFields[1].Descriptor()
	// releasepolicy.ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	releasepolicy.ModuleValidator = releasepolicyDescModule.Validators[0].(func(string) error)
	releasepolicyviolationFields := schema.ReleasePolicyViolation{}.Fields()
	_ = releasepolicyviolationFields
	// releasepolicyviolationDescMessage is the schema descriptor for message field.
	releasepolicyviolationDescMessage := releasepolicyviolationFields[0].Descriptor()
	// releasepolicyviolation.MessageValidator is a validator for the "message" field. It is called by the builders before save.
	releasepolicyviolation.MessageValidator = releasepolicyviolationDescMessage.Validators[0].(func(string) error)
	repoFields := schema.Repo{}.Fields()
	_ = repoFields
	// repoDescName is the schema descriptor for name field.
	repoDescName := repoFields[0].Descriptor()
	// repo.NameValidator is a validator for the "name" field. It is called by the builders before save.
	repo.NameValidator = repoDescName.Validators[0].(func(string) error)
	// repoDescDefaultBranch is the schema descriptor for default_branch field.
	repoDescDefaultBranch := repoFields[1].Descriptor()
	// repo.DefaultDefaultBranch holds the default value on creation for the default_branch field.
	repo.DefaultDefaultBranch = repoDescDefaultBranch.Default.(string)
	// repo.DefaultBranchValidator is a validator for the "default_branch" field. It is called by the builders before save.
	repo.DefaultBranchValidator = repoDescDefaultBranch.Validators[0].(func(string) error)
	testcaseFields := schema.TestCase{}.Fields()
	_ = testcaseFields
	// testcaseDescName is the schema descriptor for name field.
	testcaseDescName := testcaseFields[0].Descriptor()
	// testcase.NameValidator is a validator for the "name" field. It is called by the builders before save.
	testcase.NameValidator = testcaseDescName.Validators[0].(func(string) error)
	// testcaseDescMessage is the schema descriptor for message field.
	testcaseDescMessage := testcaseFields[2].Descriptor()
	// testcase.MessageValidator is a validator for the "message" field. It is called by the builders before save.
	testcase.MessageValidator = testcaseDescMessage.Validators[0].(func(string) error)
	// testcaseDescElapsed is the schema descriptor for elapsed field.
	testcaseDescElapsed := testcaseFields[3].Descriptor()
	// testcase.DefaultElapsed holds the default value on creation for the elapsed field.
	testcase.DefaultElapsed = testcaseDescElapsed.Default.(float64)
	// testcase.ElapsedValidator is a validator for the "elapsed" field. It is called by the builders before save.
	testcase.ElapsedValidator = testcaseDescElapsed.Validators[0].(func(float64) error)
	testrunHooks := schema.TestRun{}.Hooks()
	testrun.Hooks[0] = testrunHooks[0]
	testrunFields := schema.TestRun{}.Fields()
	_ = testrunFields
	// testrunDescTool is the schema descriptor for tool field.
	testrunDescTool := testrunFields[0].Descriptor()
	// testrun.ToolValidator is a validator for the "tool" field. It is called by the builders before save.
	testrun.ToolValidator = testrunDescTool.Validators[0].(func(string) error)
	// testrunDescTime is the schema descriptor for time field.
	testrunDescTime := testrunFields[1].Descriptor()
	// testrun.DefaultTime holds the default value on creation for the time field.
	testrun.DefaultTime = testrunDescTime.Default.(func() time.Time)
	vulnerabilityHooks := schema.Vulnerability{}.Hooks()
	vulnerability.Hooks[0] = vulnerabilityHooks[0]
	vulnerabilityFields := schema.Vulnerability{}.Fields()
	_ = vulnerabilityFields
	// vulnerabilityDescVid is the schema descriptor for vid field.
	vulnerabilityDescVid := vulnerabilityFields[0].Descriptor()
	// vulnerability.VidValidator is a validator for the "vid" field. It is called by the builders before save.
	vulnerability.VidValidator = vulnerabilityDescVid.Validators[0].(func(string) error)
	// vulnerabilityDescSeverityScore is the schema descriptor for severity_score field.
	vulnerabilityDescSeverityScore := vulnerabilityFields[3].Descriptor()
	// vulnerability.DefaultSeverityScore holds the default value on creation for the severity_score field.
	vulnerability.DefaultSeverityScore = vulnerabilityDescSeverityScore.Default.(float64)
	vulnerabilityreviewFields := schema.VulnerabilityReview{}.Fields()
	_ = vulnerabilityreviewFields
	// vulnerabilityreviewDescNote is the schema descriptor for note field.
	vulnerabilityreviewDescNote := vulnerabilityreviewFields[0].Descriptor()
	// vulnerabilityreview.NoteValidator is a validator for the "note" field. It is called by the builders before save.
	vulnerabilityreview.NoteValidator = vulnerabilityreviewDescNote.Validators[0].(func(string) error)
}

const (
	Version = "(devel)" // Version of ent codegen.
)
