package api

import (
	"github.com/valocode/bubbly/ent"
)

type (
	AdapterSaveRequest struct {
		Adapter *ent.AdapterModelCreate `json:"adapter,omitempty" validate:"required"`
		// *model.AdapterModel
	}

	AdapterGetRequest struct {
		Name *string `validate:"required"`
		Tag  *string
		Type *string
	}

	AdapterGetResponse struct {
		ent.AdapterModelRead `validate:"required"`
	}
)
