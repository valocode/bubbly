package api

import "github.com/valocode/bubbly/ent/model"

type (
	AdapterSaveRequest struct {
		*model.AdapterModel
	}

	AdapterGetRequest struct {
		Name *string
		Tag  *string
		Type *string
	}

	AdapterGetResponse struct {
		*model.AdapterModel
	}
)
