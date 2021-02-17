package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/events"

	"github.com/zclconf/go-cty/cty"

	"github.com/verifa/bubbly/api/common"
	"github.com/verifa/bubbly/api/core"
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

	if err := l.load(bCtx); err != nil {
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

// load creates a new client using server configurations determined from
// available viper bindings, then POSTs the load resource's Spec.Data
// to the bubbly server
// load outputs an error if any part of this process fails, nil if
// the data is successfully POSTed to the bubbly server.
func (l *Load) load(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Interface("server", bCtx.ServerConfig).Msg("loading to server with configuration")

	c, err := client.NewHTTP(bCtx)

	if err != nil {
		return fmt.Errorf("failed to establish new client for loading: %w", err)
	}

	var data core.DataBlocks

	err = json.Unmarshal([]byte(l.Spec.Data), &data)

	if err != nil {
		return fmt.Errorf("failed to unmarshal spec data for loading: %w", err)
	}

	err = c.Load(bCtx, data)

	if err != nil {
		return fmt.Errorf("failed to load spec data: %w", err)
	}

	return nil
}

type loadSpec struct {
	Inputs core.InputDeclarations `hcl:"input,block"`
	Data   string                 `hcl:"data,attr"`
}
