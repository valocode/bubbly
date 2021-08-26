package store

import (
	"fmt"

	bubblyadapter "github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/adapter"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveAdapter(req *api.AdapterSaveRequest) (*ent.Adapter, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "adapter create")
	}

	dbAdapter, err := s.client.Adapter.Query().Where(
		adapter.Name(*req.Adapter.Name), adapter.Tag(*req.Adapter.Tag),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "adapter query")
		}
		dbAdapter, err = s.client.Adapter.Create().
			SetModelCreate(req.Adapter).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "adapter create")
		}
	} else {
		dbAdapter, err = s.client.Adapter.UpdateOne(dbAdapter).
			SetModelCreate(req.Adapter).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "adater update")
		}
	}
	return dbAdapter, nil
}

func (s *Store) GetAdapter(req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, HandleValidatorError(err, "adapter read")
	}
	var tag = bubblyadapter.DefaultTag
	if req.Tag != nil {
		tag = *req.Tag
	}

	query := s.client.Adapter.Query().Where(
		adapter.Name(*req.Name), adapter.Tag(tag),
	)
	dbAdapter, err := query.Only(s.ctx)
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
