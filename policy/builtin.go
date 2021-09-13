package policy

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/valocode/bubbly/ent"
)

const policyDir = "builtin"

//go:embed builtin/*.rego
var policyFiles embed.FS

// BuiltinPolicies uses the embedded policies and returns a slice of policy
// models that can be used to populate the bubbly database
func BuiltinPolicies() ([]*ent.ReleasePolicyModelCreate, error) {
	var policies []*ent.ReleasePolicyModelCreate
	regoFiles, err := policyFiles.ReadDir(policyDir)
	if err != nil {
		return nil, fmt.Errorf("error reading builtin policies: %w", err)
	}
	for _, regoFile := range regoFiles {
		fmt.Println("rego: ", regoFile.Name())
		file, err := policyFiles.Open(filepath.Join(policyDir, regoFile.Name()))
		if err != nil {
			return nil, fmt.Errorf("opening rego file %s: %w", regoFile.Name(), err)
		}
		b, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("reading rego file %s: %w", regoFile.Name(), err)
		}
		filename := filepath.Base(regoFile.Name())
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		policy := ent.NewReleasePolicyModelCreate().
			SetName(name).
			SetModule(string(b))
		if err := Validate(*policy.Module, WithResolver(&EntResolver{})); err != nil {
			return nil, fmt.Errorf("validating policy %s: %w", regoFile.Name(), err)
		}
		policies = append(policies, policy)
	}
	return policies, nil
}
