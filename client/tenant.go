package client

import (
	"errors"
	"fmt"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/env"
)

func (h *httpClient) CreateTenant(*env.BubblyContext, *component.MessageAuth, string) error {
	return errors.New("unsupported operation for the HTTP client: CreateTenant")
}

func (n *natsClient) CreateTenant(bCtx *env.BubblyContext, auth *component.MessageAuth, tenant string) error {
	bCtx.Logger.Debug().
		Str("subject", string(component.StoreCreateTenant)).
		Msg("Posting resource to store")

	req := component.Request{
		Subject: component.StoreCreateTenant,
		Data: component.MessageData{
			Auth: auth,
			Data: []byte(tenant),
		},
	}

	if err := n.request(bCtx, &req); err != nil {
		return fmt.Errorf("failed to post resource: %w", err)
	}

	return nil
}
