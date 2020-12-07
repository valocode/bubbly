package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/env"

	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
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

func (p *Load) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *Load) Apply(bCtx *env.BubblyContext, ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(bCtx, p.SpecHCL.Body, &p.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %w`, p.String(), err),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().Msgf("Attempting to load this JSON: %s", p.Spec.Data)

	// Pull the server configuration from the resource's ResourceContext.
	// At current, this call simply creates a new config.ServerConfig
	// struct instance from a combination of defaults and viper bindings.
	_, err := bCtx.GetServerConfig()

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to establish the server to load to: %w`, err),
			Value:  cty.NilVal,
		}
	}

	err = p.load(bCtx)

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to load data to bubbly server: %w`, err),
			Value:  cty.NilVal,
		}
	}

	bCtx.Logger.Debug().Msg("JSON successfully loaded to bubbly server")

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

// load creates a new client using server configurations determined from
// available viper bindings, then POSTs the load resource's Spec.Data
// to the bubbly server
// load outputs an error if any part of this process fails, nil if
// the data is successfully POSTed to the bubbly server.
func (p *Load) load(bCtx *env.BubblyContext) error {
	bCtx.Logger.Debug().Interface("server", bCtx.ServerConfig).Msg("loading to server with configuration")

	c, err := client.New(bCtx)

	if err != nil {
		return fmt.Errorf("failed to establish new client for loading: %w", err)
	}

	var data core.DataBlocks

	err = json.Unmarshal([]byte(p.Spec.Data), &data)

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
	Inputs InputDeclarations `hcl:"input,block"`
	Data   string            `hcl:"data,attr"`
}
