package v1

import (
	"encoding/json"
	"fmt"

	"github.com/verifa/bubbly/client"
	"github.com/verifa/bubbly/config"

	"github.com/rs/zerolog/log"
	"github.com/verifa/bubbly/api/core"
	"github.com/zclconf/go-cty/cty"
)

var _ core.Publish = (*Publish)(nil)

type Publish struct {
	*core.ResourceBlock
	Spec publishSpec
}

func NewPublish(resBlock *core.ResourceBlock) *Publish {
	return &Publish{
		ResourceBlock: resBlock,
	}
}

func (p *Publish) SpecValue() core.ResourceSpec {
	return &p.Spec
}

// Apply returns ...
func (p *Publish) Apply(ctx *core.ResourceContext) core.ResourceOutput {
	if err := ctx.DecodeBody(p, p.SpecHCL.Body, &p.Spec); err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to decode "%s" body spec: %s`, p.String(), err.Error()),
			Value:  cty.NilVal,
		}
	}

	log.Debug().Msgf("Attempting to publish this JSON: %s", p.Spec.Data)

	// Pull the server configuration from the resource's ResourceContext.
	// At current, this call simply creates a new config.ServerConfig
	// struct instance from a combination of defaults and viper bindings.
	sc, err := ctx.GetServerConfig()

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to establish the server to publish to: %w`, err),
			Value:  cty.NilVal,
		}
	}

	err = p.publish(*sc)

	if err != nil {
		return core.ResourceOutput{
			Status: core.ResourceOutputFailure,
			Error:  fmt.Errorf(`Failed to publish data to bubbly server: %s`, err.Error()),
			Value:  cty.NilVal,
		}
	}

	log.Debug().Msg("JSON successfully published to bubbly server")

	return core.ResourceOutput{
		Status: core.ResourceOutputSuccess,
		Error:  nil,
		Value:  cty.NilVal,
	}
}

// publish creates a new client using server configurations determined from
// available viper bindings, then POSTs the publish resource's Spec.Data
// to the bubbly server
// publish outputs an error if any part of this process fails, nil if
// the data is successfully POSTed to the bubbly server.
func (p *Publish) publish(sc config.ServerConfig) error {
	log.Debug().Interface("cfg", sc).Msg("publishing with config")

	c, err := client.NewClient(sc)

	if err != nil {
		return err
	}

	var data core.DataBlocks

	err = json.Unmarshal([]byte(p.Spec.Data), &data)

	if err != nil {
		return err
	}

	err = c.Publish(data)

	if err != nil {
		return err
	}

	return nil
}

type publishSpec struct {
	Inputs InputDeclarations `hcl:"input,block"`
	Data   string            `hcl:"data,attr"`
}
