package bubbly

import (
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

type BubblyFileParser struct {
	Release        *releaseSpec        `hcl:"release,block"`
	ResourceBlocks core.ResourceBlocks `hcl:"resource,block"`
}

func CreateRelease(bCtx *env.BubblyContext) error {
	var fileParser BubblyFileParser
	err := parser.ParseFilename(bCtx, "./.bubbly", &fileParser)
	if err != nil {
		return err
	}

	fmt.Printf("resource: %+v\n", fileParser.Release)
	return nil
}
