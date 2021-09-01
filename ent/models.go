// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/valocode/bubbly/ent/artifact"
	"github.com/valocode/bubbly/ent/codeissue"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
	"github.com/valocode/bubbly/ent/releasepolicyviolation"
	schema "github.com/valocode/bubbly/ent/schema/types"
	"github.com/valocode/bubbly/ent/vulnerability"
)

type AdapterModelCreate struct {
	Name   *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Tag    *string `json:"tag,omitempty" validate:"required" mapstructure:"tag"`
	Module *string `json:"module,omitempty" validate:"required" mapstructure:"module"`
}

func NewAdapterModelCreate() *AdapterModelCreate {
	return &AdapterModelCreate{}
}

func (a *AdapterModelCreate) SetName(value string) *AdapterModelCreate {
	a.Name = &value
	return a
}
func (a *AdapterModelCreate) SetTag(value string) *AdapterModelCreate {
	a.Tag = &value
	return a
}
func (a *AdapterModelCreate) SetModule(value string) *AdapterModelCreate {
	a.Module = &value
	return a
}

func (a *AdapterCreate) SetModelCreate(model *AdapterModelCreate) *AdapterCreate {
	a.mutation.SetModelCreate(model)
	return a
}

func (a *AdapterUpdateOne) SetModelCreate(model *AdapterModelCreate) *AdapterUpdateOne {
	a.mutation.SetModelCreate(model)
	return a
}

func (a *AdapterMutation) SetModelCreate(model *AdapterModelCreate) *AdapterMutation {
	if model.Name != nil {
		a.SetName(*model.Name)
	}
	if model.Tag != nil {
		a.SetTag(*model.Tag)
	}
	if model.Module != nil {
		a.SetModule(*model.Module)
	}
	return a
}

type AdapterModelRead struct {
	Name   *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Tag    *string `json:"tag,omitempty" validate:"required" mapstructure:"tag"`
	Module *string `json:"module,omitempty" validate:"required" mapstructure:"module"`
	ID     *int    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewAdapterModelRead() *AdapterModelRead {
	return &AdapterModelRead{}
}

func (a *AdapterModelRead) FromEnt(value *Adapter) *AdapterModelRead {
	a.Name = &value.Name
	a.Tag = &value.Tag
	a.Module = &value.Module
	a.ID = &value.ID
	return a
}

type AdapterModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewAdapterModelUpdate() *AdapterModelUpdate {
	return &AdapterModelUpdate{}
}

func (a *AdapterModelUpdate) SetID(value int) *AdapterModelUpdate {
	a.ID = &value
	return a
}

type ArtifactModelCreate struct {
	Name     *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Sha256   *string          `json:"sha256,omitempty" validate:"required" mapstructure:"sha256"`
	Type     *artifact.Type   `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewArtifactModelCreate() *ArtifactModelCreate {
	return &ArtifactModelCreate{}
}

func (a *ArtifactModelCreate) SetName(value string) *ArtifactModelCreate {
	a.Name = &value
	return a
}
func (a *ArtifactModelCreate) SetSha256(value string) *ArtifactModelCreate {
	a.Sha256 = &value
	return a
}
func (a *ArtifactModelCreate) SetType(value artifact.Type) *ArtifactModelCreate {
	a.Type = &value
	return a
}
func (a *ArtifactModelCreate) SetTime(value time.Time) *ArtifactModelCreate {
	a.Time = &value
	return a
}
func (a *ArtifactModelCreate) SetMetadata(value schema.Metadata) *ArtifactModelCreate {
	a.Metadata = &value
	return a
}

func (a *ArtifactCreate) SetModelCreate(model *ArtifactModelCreate) *ArtifactCreate {
	a.mutation.SetModelCreate(model)
	return a
}

func (a *ArtifactUpdateOne) SetModelCreate(model *ArtifactModelCreate) *ArtifactUpdateOne {
	a.mutation.SetModelCreate(model)
	return a
}

func (a *ArtifactMutation) SetModelCreate(model *ArtifactModelCreate) *ArtifactMutation {
	if model.Name != nil {
		a.SetName(*model.Name)
	}
	if model.Sha256 != nil {
		a.SetSha256(*model.Sha256)
	}
	if model.Type != nil {
		a.SetType(*model.Type)
	}
	if model.Time != nil {
		a.SetTime(*model.Time)
	}
	if model.Metadata != nil {
		a.SetMetadata(*model.Metadata)
	}
	return a
}

type ArtifactModelRead struct {
	Name     *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Sha256   *string          `json:"sha256,omitempty" validate:"required" mapstructure:"sha256"`
	Type     *artifact.Type   `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID       *int             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewArtifactModelRead() *ArtifactModelRead {
	return &ArtifactModelRead{}
}

func (a *ArtifactModelRead) FromEnt(value *Artifact) *ArtifactModelRead {
	a.Name = &value.Name
	a.Sha256 = &value.Sha256
	a.Type = &value.Type
	a.Time = &value.Time
	a.Metadata = &value.Metadata
	a.ID = &value.ID
	return a
}

type ArtifactModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewArtifactModelUpdate() *ArtifactModelUpdate {
	return &ArtifactModelUpdate{}
}

func (a *ArtifactModelUpdate) SetID(value int) *ArtifactModelUpdate {
	a.ID = &value
	return a
}

type CodeIssueModelCreate struct {
	RuleID   *string             `json:"rule_id,omitempty" validate:"required" mapstructure:"rule_id"`
	Message  *string             `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Severity *codeissue.Severity `json:"severity,omitempty" validate:"required" mapstructure:"severity"`
	Type     *codeissue.Type     `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Metadata *schema.Metadata    `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewCodeIssueModelCreate() *CodeIssueModelCreate {
	return &CodeIssueModelCreate{}
}

func (ci *CodeIssueModelCreate) SetRuleID(value string) *CodeIssueModelCreate {
	ci.RuleID = &value
	return ci
}
func (ci *CodeIssueModelCreate) SetMessage(value string) *CodeIssueModelCreate {
	ci.Message = &value
	return ci
}
func (ci *CodeIssueModelCreate) SetSeverity(value codeissue.Severity) *CodeIssueModelCreate {
	ci.Severity = &value
	return ci
}
func (ci *CodeIssueModelCreate) SetType(value codeissue.Type) *CodeIssueModelCreate {
	ci.Type = &value
	return ci
}
func (ci *CodeIssueModelCreate) SetMetadata(value schema.Metadata) *CodeIssueModelCreate {
	ci.Metadata = &value
	return ci
}

func (ci *CodeIssueCreate) SetModelCreate(model *CodeIssueModelCreate) *CodeIssueCreate {
	ci.mutation.SetModelCreate(model)
	return ci
}

func (ci *CodeIssueUpdateOne) SetModelCreate(model *CodeIssueModelCreate) *CodeIssueUpdateOne {
	ci.mutation.SetModelCreate(model)
	return ci
}

func (ci *CodeIssueMutation) SetModelCreate(model *CodeIssueModelCreate) *CodeIssueMutation {
	if model.RuleID != nil {
		ci.SetRuleID(*model.RuleID)
	}
	if model.Message != nil {
		ci.SetMessage(*model.Message)
	}
	if model.Severity != nil {
		ci.SetSeverity(*model.Severity)
	}
	if model.Type != nil {
		ci.SetType(*model.Type)
	}
	if model.Metadata != nil {
		ci.SetMetadata(*model.Metadata)
	}
	return ci
}

type CodeIssueModelRead struct {
	RuleID   *string             `json:"rule_id,omitempty" validate:"required" mapstructure:"rule_id"`
	Message  *string             `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Severity *codeissue.Severity `json:"severity,omitempty" validate:"required" mapstructure:"severity"`
	Type     *codeissue.Type     `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Metadata *schema.Metadata    `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID       *int                `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewCodeIssueModelRead() *CodeIssueModelRead {
	return &CodeIssueModelRead{}
}

func (ci *CodeIssueModelRead) FromEnt(value *CodeIssue) *CodeIssueModelRead {
	ci.RuleID = &value.RuleID
	ci.Message = &value.Message
	ci.Severity = &value.Severity
	ci.Type = &value.Type
	ci.Metadata = &value.Metadata
	ci.ID = &value.ID
	return ci
}

type CodeIssueModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewCodeIssueModelUpdate() *CodeIssueModelUpdate {
	return &CodeIssueModelUpdate{}
}

func (ci *CodeIssueModelUpdate) SetID(value int) *CodeIssueModelUpdate {
	ci.ID = &value
	return ci
}

type CodeScanModelCreate struct {
	Tool     *string          `json:"tool,omitempty" validate:"required" mapstructure:"tool"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewCodeScanModelCreate() *CodeScanModelCreate {
	return &CodeScanModelCreate{}
}

func (cs *CodeScanModelCreate) SetTool(value string) *CodeScanModelCreate {
	cs.Tool = &value
	return cs
}
func (cs *CodeScanModelCreate) SetTime(value time.Time) *CodeScanModelCreate {
	cs.Time = &value
	return cs
}
func (cs *CodeScanModelCreate) SetMetadata(value schema.Metadata) *CodeScanModelCreate {
	cs.Metadata = &value
	return cs
}

func (cs *CodeScanCreate) SetModelCreate(model *CodeScanModelCreate) *CodeScanCreate {
	cs.mutation.SetModelCreate(model)
	return cs
}

func (cs *CodeScanUpdateOne) SetModelCreate(model *CodeScanModelCreate) *CodeScanUpdateOne {
	cs.mutation.SetModelCreate(model)
	return cs
}

func (cs *CodeScanMutation) SetModelCreate(model *CodeScanModelCreate) *CodeScanMutation {
	if model.Tool != nil {
		cs.SetTool(*model.Tool)
	}
	if model.Time != nil {
		cs.SetTime(*model.Time)
	}
	if model.Metadata != nil {
		cs.SetMetadata(*model.Metadata)
	}
	return cs
}

type CodeScanModelRead struct {
	Tool     *string          `json:"tool,omitempty" validate:"required" mapstructure:"tool"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID       *int             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewCodeScanModelRead() *CodeScanModelRead {
	return &CodeScanModelRead{}
}

func (cs *CodeScanModelRead) FromEnt(value *CodeScan) *CodeScanModelRead {
	cs.Tool = &value.Tool
	cs.Time = &value.Time
	cs.Metadata = &value.Metadata
	cs.ID = &value.ID
	return cs
}

type CodeScanModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewCodeScanModelUpdate() *CodeScanModelUpdate {
	return &CodeScanModelUpdate{}
}

func (cs *CodeScanModelUpdate) SetID(value int) *CodeScanModelUpdate {
	cs.ID = &value
	return cs
}

type ComponentModelCreate struct {
	Name        *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Vendor      *string          `json:"vendor,omitempty"  mapstructure:"vendor"`
	Version     *string          `json:"version,omitempty" validate:"required" mapstructure:"version"`
	Description *string          `json:"description,omitempty"  mapstructure:"description"`
	URL         *string          `json:"url,omitempty"  mapstructure:"url"`
	Metadata    *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewComponentModelCreate() *ComponentModelCreate {
	return &ComponentModelCreate{}
}

func (c *ComponentModelCreate) SetName(value string) *ComponentModelCreate {
	c.Name = &value
	return c
}
func (c *ComponentModelCreate) SetVendor(value string) *ComponentModelCreate {
	c.Vendor = &value
	return c
}
func (c *ComponentModelCreate) SetVersion(value string) *ComponentModelCreate {
	c.Version = &value
	return c
}
func (c *ComponentModelCreate) SetDescription(value string) *ComponentModelCreate {
	c.Description = &value
	return c
}
func (c *ComponentModelCreate) SetURL(value string) *ComponentModelCreate {
	c.URL = &value
	return c
}
func (c *ComponentModelCreate) SetMetadata(value schema.Metadata) *ComponentModelCreate {
	c.Metadata = &value
	return c
}

func (c *ComponentCreate) SetModelCreate(model *ComponentModelCreate) *ComponentCreate {
	c.mutation.SetModelCreate(model)
	return c
}

func (c *ComponentUpdateOne) SetModelCreate(model *ComponentModelCreate) *ComponentUpdateOne {
	c.mutation.SetModelCreate(model)
	return c
}

func (c *ComponentMutation) SetModelCreate(model *ComponentModelCreate) *ComponentMutation {
	if model.Name != nil {
		c.SetName(*model.Name)
	}
	if model.Vendor != nil {
		c.SetVendor(*model.Vendor)
	}
	if model.Version != nil {
		c.SetVersion(*model.Version)
	}
	if model.Description != nil {
		c.SetDescription(*model.Description)
	}
	if model.URL != nil {
		c.SetURL(*model.URL)
	}
	if model.Metadata != nil {
		c.SetMetadata(*model.Metadata)
	}
	return c
}

type ComponentModelRead struct {
	Name        *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Vendor      *string          `json:"vendor,omitempty"  mapstructure:"vendor"`
	Version     *string          `json:"version,omitempty" validate:"required" mapstructure:"version"`
	Description *string          `json:"description,omitempty"  mapstructure:"description"`
	URL         *string          `json:"url,omitempty"  mapstructure:"url"`
	Metadata    *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID          *int             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewComponentModelRead() *ComponentModelRead {
	return &ComponentModelRead{}
}

func (c *ComponentModelRead) FromEnt(value *Component) *ComponentModelRead {
	c.Name = &value.Name
	c.Vendor = &value.Vendor
	c.Version = &value.Version
	c.Description = &value.Description
	c.URL = &value.URL
	c.Metadata = &value.Metadata
	c.ID = &value.ID
	return c
}

type ComponentModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewComponentModelUpdate() *ComponentModelUpdate {
	return &ComponentModelUpdate{}
}

func (c *ComponentModelUpdate) SetID(value int) *ComponentModelUpdate {
	c.ID = &value
	return c
}

type GitCommitModelCreate struct {
	Hash   *string    `json:"hash,omitempty" validate:"required" mapstructure:"hash"`
	Branch *string    `json:"branch,omitempty" validate:"required" mapstructure:"branch"`
	Tag    *string    `json:"tag,omitempty"  mapstructure:"tag"`
	Time   *time.Time `json:"time,omitempty" validate:"required" mapstructure:"time"`
}

func NewGitCommitModelCreate() *GitCommitModelCreate {
	return &GitCommitModelCreate{}
}

func (gc *GitCommitModelCreate) SetHash(value string) *GitCommitModelCreate {
	gc.Hash = &value
	return gc
}
func (gc *GitCommitModelCreate) SetBranch(value string) *GitCommitModelCreate {
	gc.Branch = &value
	return gc
}
func (gc *GitCommitModelCreate) SetTag(value string) *GitCommitModelCreate {
	gc.Tag = &value
	return gc
}
func (gc *GitCommitModelCreate) SetTime(value time.Time) *GitCommitModelCreate {
	gc.Time = &value
	return gc
}

func (gc *GitCommitCreate) SetModelCreate(model *GitCommitModelCreate) *GitCommitCreate {
	gc.mutation.SetModelCreate(model)
	return gc
}

func (gc *GitCommitUpdateOne) SetModelCreate(model *GitCommitModelCreate) *GitCommitUpdateOne {
	gc.mutation.SetModelCreate(model)
	return gc
}

func (gc *GitCommitMutation) SetModelCreate(model *GitCommitModelCreate) *GitCommitMutation {
	if model.Hash != nil {
		gc.SetHash(*model.Hash)
	}
	if model.Branch != nil {
		gc.SetBranch(*model.Branch)
	}
	if model.Tag != nil {
		gc.SetTag(*model.Tag)
	}
	if model.Time != nil {
		gc.SetTime(*model.Time)
	}
	return gc
}

type GitCommitModelRead struct {
	Hash   *string    `json:"hash,omitempty" validate:"required" mapstructure:"hash"`
	Branch *string    `json:"branch,omitempty" validate:"required" mapstructure:"branch"`
	Tag    *string    `json:"tag,omitempty"  mapstructure:"tag"`
	Time   *time.Time `json:"time,omitempty" validate:"required" mapstructure:"time"`
	ID     *int       `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewGitCommitModelRead() *GitCommitModelRead {
	return &GitCommitModelRead{}
}

func (gc *GitCommitModelRead) FromEnt(value *GitCommit) *GitCommitModelRead {
	gc.Hash = &value.Hash
	gc.Branch = &value.Branch
	gc.Tag = &value.Tag
	gc.Time = &value.Time
	gc.ID = &value.ID
	return gc
}

type GitCommitModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewGitCommitModelUpdate() *GitCommitModelUpdate {
	return &GitCommitModelUpdate{}
}

func (gc *GitCommitModelUpdate) SetID(value int) *GitCommitModelUpdate {
	gc.ID = &value
	return gc
}

type OrganizationModelCreate struct {
	Name *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
}

func NewOrganizationModelCreate() *OrganizationModelCreate {
	return &OrganizationModelCreate{}
}

func (o *OrganizationModelCreate) SetName(value string) *OrganizationModelCreate {
	o.Name = &value
	return o
}

func (o *OrganizationCreate) SetModelCreate(model *OrganizationModelCreate) *OrganizationCreate {
	o.mutation.SetModelCreate(model)
	return o
}

func (o *OrganizationUpdateOne) SetModelCreate(model *OrganizationModelCreate) *OrganizationUpdateOne {
	o.mutation.SetModelCreate(model)
	return o
}

func (o *OrganizationMutation) SetModelCreate(model *OrganizationModelCreate) *OrganizationMutation {
	if model.Name != nil {
		o.SetName(*model.Name)
	}
	return o
}

type OrganizationModelRead struct {
	Name *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	ID   *int    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewOrganizationModelRead() *OrganizationModelRead {
	return &OrganizationModelRead{}
}

func (o *OrganizationModelRead) FromEnt(value *Organization) *OrganizationModelRead {
	o.Name = &value.Name
	o.ID = &value.ID
	return o
}

type OrganizationModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewOrganizationModelUpdate() *OrganizationModelUpdate {
	return &OrganizationModelUpdate{}
}

func (o *OrganizationModelUpdate) SetID(value int) *OrganizationModelUpdate {
	o.ID = &value
	return o
}

type ProjectModelCreate struct {
	Name *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
}

func NewProjectModelCreate() *ProjectModelCreate {
	return &ProjectModelCreate{}
}

func (pr *ProjectModelCreate) SetName(value string) *ProjectModelCreate {
	pr.Name = &value
	return pr
}

func (pr *ProjectCreate) SetModelCreate(model *ProjectModelCreate) *ProjectCreate {
	pr.mutation.SetModelCreate(model)
	return pr
}

func (pr *ProjectUpdateOne) SetModelCreate(model *ProjectModelCreate) *ProjectUpdateOne {
	pr.mutation.SetModelCreate(model)
	return pr
}

func (pr *ProjectMutation) SetModelCreate(model *ProjectModelCreate) *ProjectMutation {
	if model.Name != nil {
		pr.SetName(*model.Name)
	}
	return pr
}

type ProjectModelRead struct {
	Name *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	ID   *int    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewProjectModelRead() *ProjectModelRead {
	return &ProjectModelRead{}
}

func (pr *ProjectModelRead) FromEnt(value *Project) *ProjectModelRead {
	pr.Name = &value.Name
	pr.ID = &value.ID
	return pr
}

type ProjectModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewProjectModelUpdate() *ProjectModelUpdate {
	return &ProjectModelUpdate{}
}

func (pr *ProjectModelUpdate) SetID(value int) *ProjectModelUpdate {
	pr.ID = &value
	return pr
}

type ReleaseModelCreate struct {
	Name    *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Version *string `json:"version,omitempty" validate:"required" mapstructure:"version"`
}

func NewReleaseModelCreate() *ReleaseModelCreate {
	return &ReleaseModelCreate{}
}

func (r *ReleaseModelCreate) SetName(value string) *ReleaseModelCreate {
	r.Name = &value
	return r
}
func (r *ReleaseModelCreate) SetVersion(value string) *ReleaseModelCreate {
	r.Version = &value
	return r
}

func (r *ReleaseCreate) SetModelCreate(model *ReleaseModelCreate) *ReleaseCreate {
	r.mutation.SetModelCreate(model)
	return r
}

func (r *ReleaseUpdateOne) SetModelCreate(model *ReleaseModelCreate) *ReleaseUpdateOne {
	r.mutation.SetModelCreate(model)
	return r
}

func (r *ReleaseMutation) SetModelCreate(model *ReleaseModelCreate) *ReleaseMutation {
	if model.Name != nil {
		r.SetName(*model.Name)
	}
	if model.Version != nil {
		r.SetVersion(*model.Version)
	}
	return r
}

type ReleaseModelRead struct {
	Name    *string         `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Version *string         `json:"version,omitempty" validate:"required" mapstructure:"version"`
	Status  *release.Status `json:"status,omitempty"  mapstructure:"status"`
	ID      *int            `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleaseModelRead() *ReleaseModelRead {
	return &ReleaseModelRead{}
}

func (r *ReleaseModelRead) FromEnt(value *Release) *ReleaseModelRead {
	r.Name = &value.Name
	r.Version = &value.Version
	r.Status = &value.Status
	r.ID = &value.ID
	return r
}

type ReleaseModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleaseModelUpdate() *ReleaseModelUpdate {
	return &ReleaseModelUpdate{}
}

func (r *ReleaseModelUpdate) SetID(value int) *ReleaseModelUpdate {
	r.ID = &value
	return r
}

type ReleaseEntryModelCreate struct {
	Type *releaseentry.Type `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Time *time.Time         `json:"time,omitempty"  mapstructure:"time"`
}

func NewReleaseEntryModelCreate() *ReleaseEntryModelCreate {
	return &ReleaseEntryModelCreate{}
}

func (re *ReleaseEntryModelCreate) SetType(value releaseentry.Type) *ReleaseEntryModelCreate {
	re.Type = &value
	return re
}
func (re *ReleaseEntryModelCreate) SetTime(value time.Time) *ReleaseEntryModelCreate {
	re.Time = &value
	return re
}

func (re *ReleaseEntryCreate) SetModelCreate(model *ReleaseEntryModelCreate) *ReleaseEntryCreate {
	re.mutation.SetModelCreate(model)
	return re
}

func (re *ReleaseEntryUpdateOne) SetModelCreate(model *ReleaseEntryModelCreate) *ReleaseEntryUpdateOne {
	re.mutation.SetModelCreate(model)
	return re
}

func (re *ReleaseEntryMutation) SetModelCreate(model *ReleaseEntryModelCreate) *ReleaseEntryMutation {
	if model.Type != nil {
		re.SetType(*model.Type)
	}
	if model.Time != nil {
		re.SetTime(*model.Time)
	}
	return re
}

type ReleaseEntryModelRead struct {
	Type *releaseentry.Type `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Time *time.Time         `json:"time,omitempty"  mapstructure:"time"`
	ID   *int               `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleaseEntryModelRead() *ReleaseEntryModelRead {
	return &ReleaseEntryModelRead{}
}

func (re *ReleaseEntryModelRead) FromEnt(value *ReleaseEntry) *ReleaseEntryModelRead {
	re.Type = &value.Type
	re.Time = &value.Time
	re.ID = &value.ID
	return re
}

type ReleaseEntryModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleaseEntryModelUpdate() *ReleaseEntryModelUpdate {
	return &ReleaseEntryModelUpdate{}
}

func (re *ReleaseEntryModelUpdate) SetID(value int) *ReleaseEntryModelUpdate {
	re.ID = &value
	return re
}

type ReleasePolicyModelCreate struct {
	Name   *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Module *string `json:"module,omitempty" validate:"required" mapstructure:"module"`
}

func NewReleasePolicyModelCreate() *ReleasePolicyModelCreate {
	return &ReleasePolicyModelCreate{}
}

func (rp *ReleasePolicyModelCreate) SetName(value string) *ReleasePolicyModelCreate {
	rp.Name = &value
	return rp
}
func (rp *ReleasePolicyModelCreate) SetModule(value string) *ReleasePolicyModelCreate {
	rp.Module = &value
	return rp
}

func (rp *ReleasePolicyCreate) SetModelCreate(model *ReleasePolicyModelCreate) *ReleasePolicyCreate {
	rp.mutation.SetModelCreate(model)
	return rp
}

func (rp *ReleasePolicyUpdateOne) SetModelCreate(model *ReleasePolicyModelCreate) *ReleasePolicyUpdateOne {
	rp.mutation.SetModelCreate(model)
	return rp
}

func (rp *ReleasePolicyMutation) SetModelCreate(model *ReleasePolicyModelCreate) *ReleasePolicyMutation {
	if model.Name != nil {
		rp.SetName(*model.Name)
	}
	if model.Module != nil {
		rp.SetModule(*model.Module)
	}
	return rp
}

type ReleasePolicyModelRead struct {
	Name   *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Module *string `json:"module,omitempty" validate:"required" mapstructure:"module"`
	ID     *int    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleasePolicyModelRead() *ReleasePolicyModelRead {
	return &ReleasePolicyModelRead{}
}

func (rp *ReleasePolicyModelRead) FromEnt(value *ReleasePolicy) *ReleasePolicyModelRead {
	rp.Name = &value.Name
	rp.Module = &value.Module
	rp.ID = &value.ID
	return rp
}

type ReleasePolicyModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleasePolicyModelUpdate() *ReleasePolicyModelUpdate {
	return &ReleasePolicyModelUpdate{}
}

func (rp *ReleasePolicyModelUpdate) SetID(value int) *ReleasePolicyModelUpdate {
	rp.ID = &value
	return rp
}

type ReleasePolicyViolationModelCreate struct {
	Message  *string                          `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Type     *releasepolicyviolation.Type     `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Severity *releasepolicyviolation.Severity `json:"severity,omitempty" validate:"required" mapstructure:"severity"`
}

func NewReleasePolicyViolationModelCreate() *ReleasePolicyViolationModelCreate {
	return &ReleasePolicyViolationModelCreate{}
}

func (rpv *ReleasePolicyViolationModelCreate) SetMessage(value string) *ReleasePolicyViolationModelCreate {
	rpv.Message = &value
	return rpv
}
func (rpv *ReleasePolicyViolationModelCreate) SetType(value releasepolicyviolation.Type) *ReleasePolicyViolationModelCreate {
	rpv.Type = &value
	return rpv
}
func (rpv *ReleasePolicyViolationModelCreate) SetSeverity(value releasepolicyviolation.Severity) *ReleasePolicyViolationModelCreate {
	rpv.Severity = &value
	return rpv
}

func (rpv *ReleasePolicyViolationCreate) SetModelCreate(model *ReleasePolicyViolationModelCreate) *ReleasePolicyViolationCreate {
	rpv.mutation.SetModelCreate(model)
	return rpv
}

func (rpv *ReleasePolicyViolationUpdateOne) SetModelCreate(model *ReleasePolicyViolationModelCreate) *ReleasePolicyViolationUpdateOne {
	rpv.mutation.SetModelCreate(model)
	return rpv
}

func (rpv *ReleasePolicyViolationMutation) SetModelCreate(model *ReleasePolicyViolationModelCreate) *ReleasePolicyViolationMutation {
	if model.Message != nil {
		rpv.SetMessage(*model.Message)
	}
	if model.Type != nil {
		rpv.SetType(*model.Type)
	}
	if model.Severity != nil {
		rpv.SetSeverity(*model.Severity)
	}
	return rpv
}

type ReleasePolicyViolationModelRead struct {
	Message  *string                          `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Type     *releasepolicyviolation.Type     `json:"type,omitempty" validate:"required" mapstructure:"type"`
	Severity *releasepolicyviolation.Severity `json:"severity,omitempty" validate:"required" mapstructure:"severity"`
	ID       *int                             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleasePolicyViolationModelRead() *ReleasePolicyViolationModelRead {
	return &ReleasePolicyViolationModelRead{}
}

func (rpv *ReleasePolicyViolationModelRead) FromEnt(value *ReleasePolicyViolation) *ReleasePolicyViolationModelRead {
	rpv.Message = &value.Message
	rpv.Type = &value.Type
	rpv.Severity = &value.Severity
	rpv.ID = &value.ID
	return rpv
}

type ReleasePolicyViolationModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewReleasePolicyViolationModelUpdate() *ReleasePolicyViolationModelUpdate {
	return &ReleasePolicyViolationModelUpdate{}
}

func (rpv *ReleasePolicyViolationModelUpdate) SetID(value int) *ReleasePolicyViolationModelUpdate {
	rpv.ID = &value
	return rpv
}

type RepoModelCreate struct {
	Name          *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	DefaultBranch *string `json:"default_branch,omitempty"  mapstructure:"default_branch"`
}

func NewRepoModelCreate() *RepoModelCreate {
	return &RepoModelCreate{}
}

func (r *RepoModelCreate) SetName(value string) *RepoModelCreate {
	r.Name = &value
	return r
}
func (r *RepoModelCreate) SetDefaultBranch(value string) *RepoModelCreate {
	r.DefaultBranch = &value
	return r
}

func (r *RepoCreate) SetModelCreate(model *RepoModelCreate) *RepoCreate {
	r.mutation.SetModelCreate(model)
	return r
}

func (r *RepoUpdateOne) SetModelCreate(model *RepoModelCreate) *RepoUpdateOne {
	r.mutation.SetModelCreate(model)
	return r
}

func (r *RepoMutation) SetModelCreate(model *RepoModelCreate) *RepoMutation {
	if model.Name != nil {
		r.SetName(*model.Name)
	}
	if model.DefaultBranch != nil {
		r.SetDefaultBranch(*model.DefaultBranch)
	}
	return r
}

type RepoModelRead struct {
	Name          *string `json:"name,omitempty" validate:"required" mapstructure:"name"`
	DefaultBranch *string `json:"default_branch,omitempty"  mapstructure:"default_branch"`
	ID            *int    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewRepoModelRead() *RepoModelRead {
	return &RepoModelRead{}
}

func (r *RepoModelRead) FromEnt(value *Repo) *RepoModelRead {
	r.Name = &value.Name
	r.DefaultBranch = &value.DefaultBranch
	r.ID = &value.ID
	return r
}

type RepoModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewRepoModelUpdate() *RepoModelUpdate {
	return &RepoModelUpdate{}
}

func (r *RepoModelUpdate) SetID(value int) *RepoModelUpdate {
	r.ID = &value
	return r
}

type TestCaseModelCreate struct {
	Name     *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Result   *bool            `json:"result,omitempty" validate:"required" mapstructure:"result"`
	Message  *string          `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Elapsed  *float64         `json:"elapsed,omitempty"  mapstructure:"elapsed"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewTestCaseModelCreate() *TestCaseModelCreate {
	return &TestCaseModelCreate{}
}

func (tc *TestCaseModelCreate) SetName(value string) *TestCaseModelCreate {
	tc.Name = &value
	return tc
}
func (tc *TestCaseModelCreate) SetResult(value bool) *TestCaseModelCreate {
	tc.Result = &value
	return tc
}
func (tc *TestCaseModelCreate) SetMessage(value string) *TestCaseModelCreate {
	tc.Message = &value
	return tc
}
func (tc *TestCaseModelCreate) SetElapsed(value float64) *TestCaseModelCreate {
	tc.Elapsed = &value
	return tc
}
func (tc *TestCaseModelCreate) SetMetadata(value schema.Metadata) *TestCaseModelCreate {
	tc.Metadata = &value
	return tc
}

func (tc *TestCaseCreate) SetModelCreate(model *TestCaseModelCreate) *TestCaseCreate {
	tc.mutation.SetModelCreate(model)
	return tc
}

func (tc *TestCaseUpdateOne) SetModelCreate(model *TestCaseModelCreate) *TestCaseUpdateOne {
	tc.mutation.SetModelCreate(model)
	return tc
}

func (tc *TestCaseMutation) SetModelCreate(model *TestCaseModelCreate) *TestCaseMutation {
	if model.Name != nil {
		tc.SetName(*model.Name)
	}
	if model.Result != nil {
		tc.SetResult(*model.Result)
	}
	if model.Message != nil {
		tc.SetMessage(*model.Message)
	}
	if model.Elapsed != nil {
		tc.SetElapsed(*model.Elapsed)
	}
	if model.Metadata != nil {
		tc.SetMetadata(*model.Metadata)
	}
	return tc
}

type TestCaseModelRead struct {
	Name     *string          `json:"name,omitempty" validate:"required" mapstructure:"name"`
	Result   *bool            `json:"result,omitempty" validate:"required" mapstructure:"result"`
	Message  *string          `json:"message,omitempty" validate:"required" mapstructure:"message"`
	Elapsed  *float64         `json:"elapsed,omitempty"  mapstructure:"elapsed"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID       *int             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewTestCaseModelRead() *TestCaseModelRead {
	return &TestCaseModelRead{}
}

func (tc *TestCaseModelRead) FromEnt(value *TestCase) *TestCaseModelRead {
	tc.Name = &value.Name
	tc.Result = &value.Result
	tc.Message = &value.Message
	tc.Elapsed = &value.Elapsed
	tc.Metadata = &value.Metadata
	tc.ID = &value.ID
	return tc
}

type TestCaseModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewTestCaseModelUpdate() *TestCaseModelUpdate {
	return &TestCaseModelUpdate{}
}

func (tc *TestCaseModelUpdate) SetID(value int) *TestCaseModelUpdate {
	tc.ID = &value
	return tc
}

type TestRunModelCreate struct {
	Tool     *string          `json:"tool,omitempty" validate:"required" mapstructure:"tool"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewTestRunModelCreate() *TestRunModelCreate {
	return &TestRunModelCreate{}
}

func (tr *TestRunModelCreate) SetTool(value string) *TestRunModelCreate {
	tr.Tool = &value
	return tr
}
func (tr *TestRunModelCreate) SetTime(value time.Time) *TestRunModelCreate {
	tr.Time = &value
	return tr
}
func (tr *TestRunModelCreate) SetMetadata(value schema.Metadata) *TestRunModelCreate {
	tr.Metadata = &value
	return tr
}

func (tr *TestRunCreate) SetModelCreate(model *TestRunModelCreate) *TestRunCreate {
	tr.mutation.SetModelCreate(model)
	return tr
}

func (tr *TestRunUpdateOne) SetModelCreate(model *TestRunModelCreate) *TestRunUpdateOne {
	tr.mutation.SetModelCreate(model)
	return tr
}

func (tr *TestRunMutation) SetModelCreate(model *TestRunModelCreate) *TestRunMutation {
	if model.Tool != nil {
		tr.SetTool(*model.Tool)
	}
	if model.Time != nil {
		tr.SetTime(*model.Time)
	}
	if model.Metadata != nil {
		tr.SetMetadata(*model.Metadata)
	}
	return tr
}

type TestRunModelRead struct {
	Tool     *string          `json:"tool,omitempty" validate:"required" mapstructure:"tool"`
	Time     *time.Time       `json:"time,omitempty"  mapstructure:"time"`
	Metadata *schema.Metadata `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID       *int             `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewTestRunModelRead() *TestRunModelRead {
	return &TestRunModelRead{}
}

func (tr *TestRunModelRead) FromEnt(value *TestRun) *TestRunModelRead {
	tr.Tool = &value.Tool
	tr.Time = &value.Time
	tr.Metadata = &value.Metadata
	tr.ID = &value.ID
	return tr
}

type TestRunModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewTestRunModelUpdate() *TestRunModelUpdate {
	return &TestRunModelUpdate{}
}

func (tr *TestRunModelUpdate) SetID(value int) *TestRunModelUpdate {
	tr.ID = &value
	return tr
}

type VulnerabilityModelCreate struct {
	Vid           *string                 `json:"vid,omitempty" validate:"required" mapstructure:"vid"`
	Summary       *string                 `json:"summary,omitempty"  mapstructure:"summary"`
	Description   *string                 `json:"description,omitempty"  mapstructure:"description"`
	SeverityScore *float64                `json:"severity_score,omitempty"  mapstructure:"severity_score"`
	Severity      *vulnerability.Severity `json:"severity,omitempty"  mapstructure:"severity"`
	Published     *time.Time              `json:"published,omitempty"  mapstructure:"published"`
	Modified      *time.Time              `json:"modified,omitempty"  mapstructure:"modified"`
	Metadata      *schema.Metadata        `json:"metadata,omitempty"  mapstructure:"metadata"`
}

func NewVulnerabilityModelCreate() *VulnerabilityModelCreate {
	return &VulnerabilityModelCreate{}
}

func (v *VulnerabilityModelCreate) SetVid(value string) *VulnerabilityModelCreate {
	v.Vid = &value
	return v
}
func (v *VulnerabilityModelCreate) SetSummary(value string) *VulnerabilityModelCreate {
	v.Summary = &value
	return v
}
func (v *VulnerabilityModelCreate) SetDescription(value string) *VulnerabilityModelCreate {
	v.Description = &value
	return v
}
func (v *VulnerabilityModelCreate) SetSeverityScore(value float64) *VulnerabilityModelCreate {
	v.SeverityScore = &value
	return v
}
func (v *VulnerabilityModelCreate) SetSeverity(value vulnerability.Severity) *VulnerabilityModelCreate {
	v.Severity = &value
	return v
}
func (v *VulnerabilityModelCreate) SetPublished(value time.Time) *VulnerabilityModelCreate {
	v.Published = &value
	return v
}
func (v *VulnerabilityModelCreate) SetModified(value time.Time) *VulnerabilityModelCreate {
	v.Modified = &value
	return v
}
func (v *VulnerabilityModelCreate) SetMetadata(value schema.Metadata) *VulnerabilityModelCreate {
	v.Metadata = &value
	return v
}

func (v *VulnerabilityCreate) SetModelCreate(model *VulnerabilityModelCreate) *VulnerabilityCreate {
	v.mutation.SetModelCreate(model)
	return v
}

func (v *VulnerabilityUpdateOne) SetModelCreate(model *VulnerabilityModelCreate) *VulnerabilityUpdateOne {
	v.mutation.SetModelCreate(model)
	return v
}

func (v *VulnerabilityMutation) SetModelCreate(model *VulnerabilityModelCreate) *VulnerabilityMutation {
	if model.Vid != nil {
		v.SetVid(*model.Vid)
	}
	if model.Summary != nil {
		v.SetSummary(*model.Summary)
	}
	if model.Description != nil {
		v.SetDescription(*model.Description)
	}
	if model.SeverityScore != nil {
		v.SetSeverityScore(*model.SeverityScore)
	}
	if model.Severity != nil {
		v.SetSeverity(*model.Severity)
	}
	if model.Published != nil {
		v.SetPublished(*model.Published)
	}
	if model.Modified != nil {
		v.SetModified(*model.Modified)
	}
	if model.Metadata != nil {
		v.SetMetadata(*model.Metadata)
	}
	return v
}

type VulnerabilityModelRead struct {
	Vid           *string                 `json:"vid,omitempty" validate:"required" mapstructure:"vid"`
	Summary       *string                 `json:"summary,omitempty"  mapstructure:"summary"`
	Description   *string                 `json:"description,omitempty"  mapstructure:"description"`
	SeverityScore *float64                `json:"severity_score,omitempty"  mapstructure:"severity_score"`
	Severity      *vulnerability.Severity `json:"severity,omitempty"  mapstructure:"severity"`
	Published     *time.Time              `json:"published,omitempty"  mapstructure:"published"`
	Modified      *time.Time              `json:"modified,omitempty"  mapstructure:"modified"`
	Metadata      *schema.Metadata        `json:"metadata,omitempty"  mapstructure:"metadata"`
	ID            *int                    `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewVulnerabilityModelRead() *VulnerabilityModelRead {
	return &VulnerabilityModelRead{}
}

func (v *VulnerabilityModelRead) FromEnt(value *Vulnerability) *VulnerabilityModelRead {
	v.Vid = &value.Vid
	v.Summary = &value.Summary
	v.Description = &value.Description
	v.SeverityScore = &value.SeverityScore
	v.Severity = &value.Severity
	v.Published = &value.Published
	v.Modified = &value.Modified
	v.Metadata = &value.Metadata
	v.ID = &value.ID
	return v
}

type VulnerabilityModelUpdate struct {
	ID *int `json:"id,omitempty" validate:"required" mapstructure:"id"`
}

func NewVulnerabilityModelUpdate() *VulnerabilityModelUpdate {
	return &VulnerabilityModelUpdate{}
}

func (v *VulnerabilityModelUpdate) SetID(value int) *VulnerabilityModelUpdate {
	v.ID = &value
	return v
}
