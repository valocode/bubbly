package v1

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"

	"github.com/zclconf/go-cty/cty"

	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
)

var _ core.Load = (*Load)(nil)

type Load struct {
	*core.ResourceBlock
	Spec loadSpec
}

func NewLoad(resBlock *core.ResourceBlock) *Load {
	return &Load{
		ResourceBlock: resBlock,
	}
}

func (l *Load) SpecValue() core.ResourceSpec {
	return &l.Spec
}

// Apply returns ...
func (l *Load) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := common.DecodeBodyWithInputs(bCtx, l.SpecHCL.Body, &l.Spec, ctx); err != nil {
		return core.ResourceOutput{
			ID:     l.String(),
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to decode "%s" body spec: %s`, l.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}
	var data core.DataBlocks
	if err := json.Unmarshal([]byte(l.Spec.Data), &data); err != nil {
		return core.ResourceOutput{
			ID:     l.String(),
			Status: events.ResourceRunFailure,
			Error:  fmt.Errorf("error unmarshalling data to data blocks: %w", err),
			Value:  cty.NilVal,
		}
	}

	// If there is contextual data that should be loaded, prefix this to the data
	// so that any joins to this data will succeed
	if ctx.DataBlocks != nil {
		data = append(ctx.DataBlocks, data...)
	}

	if err := l.load(bCtx, ctx, data); err != nil {
		return core.ResourceOutput{
			ID:     l.String(),
			Status: events.ResourceApplyFailure,
			Error:  fmt.Errorf(`failed to load data to bubbly server: %w`, err),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().Msg("JSON successfully loaded to bubbly server")

	return core.ResourceOutput{
		ID:     l.String(),
		Status: events.ResourceApplySuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

// load creates a new client then POSTs the load resource's Spec.Data
// to the bubbly server.
// load outputs an error if any part of this process fails, nil if
// the data is successfully POSTed to the bubbly server.
func (l *Load) load(bCtx *env.BubblyContext, ctx *core.ResourceContext, data core.DataBlocks) error {
	bCtx.Logger.Debug().Interface("server", bCtx.ServerConfig).Msg("loading to server with configuration")

	c, err := client.New(bCtx)
	if err != nil {
		return fmt.Errorf("failed to establish new client for loading: %w", err)
	}
	defer c.Close()

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data blocks: %w", err)
	}

	err = c.Load(bCtx, ctx.Auth, bytes)
	if err != nil {
		return fmt.Errorf("failed to load spec data: %w", err)
	}

	return nil
}

type loadSpec struct {
	Inputs core.InputDeclarations `hcl:"input,block"`
	Data   string                 `hcl:"data,attr"`
	// GitItem gitItem                `hcl:"git,block"`
}
