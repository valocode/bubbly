package store

import (
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/adapter"
	"github.com/valocode/bubbly/store/api"
)

type AdapterQuery struct {
	Where *ent.AdapterWhereInput
}

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
			SetOwnerID(h.orgID).
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

func (h *Handler) GetAdapters(query *AdapterQuery) ([]*api.Adapter, error) {
	dbAdapters, err := h.client.Adapter.Query().WhereInput(query.Where).All(h.ctx)
	if err != nil {
		return nil, HandleEntError(err, "get adapters")
	}
	adapters := make([]*api.Adapter, 0, len(dbAdapters))
	for _, a := range dbAdapters {
		adapters = append(adapters, &api.Adapter{
			AdapterModelRead: *ent.NewAdapterModelRead().FromEnt(a),
		})
	}
	return adapters, nil
}
