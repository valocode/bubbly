package adapter

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type (
	Operation struct {
		Type string        `json:"type,omitempty" hcl:",label"`
		Body hcl.Body      `json:"-" hcl:",remain"`
		Spec OperationSpec `json:"spec,omitempty"`
	}

	OperationSpec interface {
		Run(name string, args RunArgs) (cty.Value, error)
		Decode(body hcl.Body) error
	}
)

func (o *Operation) Decode() error {
	switch OperationType(o.Type) {
	case OpJSON:
		o.Spec = &jsonOp{}
	default:
		return fmt.Errorf("unsupported operation type \"%s\"", o.Type)
	}

	if err := o.Spec.Decode(o.Body); err != nil {
		return fmt.Errorf("decoding operation: %w", err)
	}
	return nil
}

func (o *Operation) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &jsonMap)
	if err != nil {
		return err
	}
	switch OperationType(o.Type) {
	case OpJSON:
		o.Spec = &jsonOp{}
	default:
		return errors.New("must set type of operation before unmarshalling")
	}
	if js, ok := jsonMap["spec"]; ok && js != nil {
		err := json.Unmarshal(*js, o.Spec)
		if err != nil {
			return err
		}
	}
	return nil
}

type OperationType string

const (
	OpHTTP OperationType = "http"
	OpJSON OperationType = "json"
	OpXML  OperationType = "xml"
)

func (a OperationType) String() string {
	return string(a)
}

func NewEvalContext(variables map[string]cty.Value) *hcl.EvalContext {
	if variables == nil {
		variables = make(map[string]cty.Value)
	}
	return &hcl.EvalContext{
		Variables: variables,
		Functions: stdfunctions(),
	}
}
