package main

import "github.com/zclconf/go-cty/cty"

type Resource interface {
	Decode() error
	Output() ResourceOutput
}

type ResourceOutput interface {
	CtyValue() cty.Value
}
