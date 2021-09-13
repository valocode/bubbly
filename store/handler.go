package store

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/organization"
)

type handlerConfig struct {
	ctx       context.Context
	client    *ent.Client
	validator *validator.Validate
	orgName   string
	userID    string
}

type Handler struct {
	ctx       context.Context
	client    *ent.Client
	validator *validator.Validate
	orgID     int
	userID    string
}

func NewHandler(opts ...func(h *handlerConfig)) (*Handler, error) {
	conf := handlerConfig{}
	for _, opt := range opts {
		opt(&conf)
	}

	h := Handler{
		ctx:       conf.ctx,
		client:    conf.client,
		validator: conf.validator,
		userID:    conf.userID,
	}
	if h.ctx == nil {
		h.ctx = context.Background()
	}
	if h.client == nil {
		return nil, errors.New("store handler requires an ent client")
	}
	if h.validator == nil {
		h.validator = newValidator()
	}
	if h.userID == "" {
		// If no user id was provided, setup annonymous user access, and check
		// that it has permissions to perform the request
		// TODO: h.userID = "some annonymous user"
	}

	orgName := conf.orgName
	if orgName == "" {
		orgName = config.DefaultOrganization
	}
	// Check that the organization exists, because all operations have to happen within an organization
	orgID, err := h.client.Organization.Query().
		Where(organization.Name(orgName)).
		OnlyID(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "fetching organization")
	}
	h.orgID = orgID

	return &h, nil
}

func (h *Handler) Client() *ent.Client {
	return h.client
}

func WithStore(s *Store) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.ctx = s.ctx
		h.client = s.client
		h.validator = s.validator
	}
}

func WithContext(ctx context.Context) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.ctx = ctx
	}
}

func WithClient(client *ent.Client) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.client = client
	}
}

func WithOrgName(name string) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.orgName = name
	}
}

func WithUserID(userID string) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.userID = userID
	}
}

func WithValidator(validator *validator.Validate) func(h *handlerConfig) {
	return func(h *handlerConfig) {
		h.validator = validator
	}
}
