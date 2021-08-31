package store

import (
	"fmt"

	bubblyadapter "github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/adapter"
	"github.com/valocode/bubbly/store/api"
)

func (h *Handler) SaveAdapter(req *api.AdapterSaveRequest) (*ent.Adapter, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "adapter create")
	}

	dbAdapter, err := h.client.Adapter.Query().Where(
		adapter.Name(*req.Adapter.Name), adapter.Tag(*req.Adapter.Tag),
	).Only(h.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "adapter query")
		}
		dbAdapter, err = h.client.Adapter.Create().
			SetModelCreate(req.Adapter).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "adapter create")
		}
	} else {
		dbAdapter, err = h.client.Adapter.UpdateOne(dbAdapter).
			SetModelCreate(req.Adapter).
			Save(h.ctx)
		if err != nil {
			return nil, HandleEntError(err, "adater update")
		}
	}
	return dbAdapter, nil
}

func (h *Handler) GetAdapter(req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "adapter read")
	}
	var tag = bubblyadapter.DefaultTag
	if req.Tag != nil {
		tag = *req.Tag
	}

	query := h.client.Adapter.Query().Where(
		adapter.Name(*req.Name), adapter.Tag(tag),
	)
	dbAdapter, err := query.Only(h.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError(err, "adapter")
		}
		return nil, HandleEntError(err, "adapter")
	}
	fmt.Println("got adapter: " + dbAdapter.String())
	return &api.AdapterGetResponse{
		AdapterModelRead: *ent.NewAdapterModelRead().FromEnt(dbAdapter),
	}, nil
}
