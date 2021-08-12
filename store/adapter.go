package store

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	bubblyadapter "github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/adapter"
	"github.com/valocode/bubbly/ent/model"
	"github.com/valocode/bubbly/store/api"
)

func (s *Store) SaveAdapter(req *api.AdapterSaveRequest) (*ent.Adapter, error) {
	var (
		vErr multierror.Error
		tag  = bubblyadapter.DefaultTag
	)
	if req.Name == nil {
		vErr = *multierror.Append(&vErr, errors.New("name required"))
	}
	if req.Type == nil {
		vErr = *multierror.Append(&vErr, errors.New("type required"))
	}
	if req.Operation == nil {
		vErr = *multierror.Append(&vErr, errors.New("operation required"))
	}
	if req.ResultsType == nil {
		vErr = *multierror.Append(&vErr, errors.New("results_type required"))
	}
	if req.Results == nil {
		vErr = *multierror.Append(&vErr, errors.New("results required"))
	}
	if vErr.ErrorOrNil() != nil {
		return nil, HandleMultiVError(vErr)
	}

	// TODO: add go playground validator...
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.Tag != nil && *req.Tag != "" {
		tag = *req.Tag
	}
	dbAdapter, err := s.client.Adapter.Query().Where(
		adapter.Name(*req.Name), adapter.Tag(tag),
	).Only(s.ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, HandleEntError(err, "adapter")
		}

		// TODO: CHECK ALL THE VALUES FOR NIL
		fmt.Printf("Results: %s\n", *req.Results)
		dbAdapter, err = s.client.Adapter.Create().
			SetName(*req.Name).
			SetTag(tag).
			SetType(*req.Type).
			SetOperation(*req.Operation).
			SetResults(*req.Results).
			SetResultsType(*req.ResultsType).
			Save(s.ctx)
		if err != nil {
			return nil, HandleEntError(err, "adapter")
		}
	}
	return dbAdapter, nil
}

func (s *Store) GetAdapter(req *api.AdapterGetRequest) (*api.AdapterGetResponse, error) {
	var tag = bubblyadapter.DefaultTag
	if req.Name == nil {
		return nil, NewValidationError(nil, "name required")
	}
	if req.Tag != nil {
		tag = *req.Tag
	}

	query := s.client.Adapter.Query().Where(
		adapter.Name(*req.Name), adapter.Tag(tag),
	)
	if req.Type != nil {
		query.Where(adapter.TypeEQ(adapter.Type(*req.Type)))
	}
	dbAdapter, err := query.Only(s.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, NewNotFoundError(err, "adapter")
		}
		return nil, HandleEntError(err, "adapter")
	}
	fmt.Println("got adapter: " + dbAdapter.String())
	return &api.AdapterGetResponse{
		AdapterModel: model.NewAdapterModel().
			SetName(dbAdapter.Name).
			SetTag(dbAdapter.Tag).
			SetType(dbAdapter.Type).
			SetOperation(dbAdapter.Operation).
			SetResultsType(dbAdapter.ResultsType).
			SetResults(dbAdapter.Results),
	}, nil
}
