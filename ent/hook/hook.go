// Code generated by entc, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/valocode/bubbly/ent"
)

// The ArtifactFunc type is an adapter to allow the use of ordinary
// function as Artifact mutator.
type ArtifactFunc func(context.Context, *ent.ArtifactMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ArtifactFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ArtifactMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ArtifactMutation", m)
	}
	return f(ctx, mv)
}

// The CVEFunc type is an adapter to allow the use of ordinary
// function as CVE mutator.
type CVEFunc func(context.Context, *ent.CVEMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CVEFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CVEMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CVEMutation", m)
	}
	return f(ctx, mv)
}

// The CVERuleFunc type is an adapter to allow the use of ordinary
// function as CVERule mutator.
type CVERuleFunc func(context.Context, *ent.CVERuleMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CVERuleFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CVERuleMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CVERuleMutation", m)
	}
	return f(ctx, mv)
}

// The CWEFunc type is an adapter to allow the use of ordinary
// function as CWE mutator.
type CWEFunc func(context.Context, *ent.CWEMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CWEFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CWEMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CWEMutation", m)
	}
	return f(ctx, mv)
}

// The CodeIssueFunc type is an adapter to allow the use of ordinary
// function as CodeIssue mutator.
type CodeIssueFunc func(context.Context, *ent.CodeIssueMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CodeIssueFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CodeIssueMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CodeIssueMutation", m)
	}
	return f(ctx, mv)
}

// The CodeScanFunc type is an adapter to allow the use of ordinary
// function as CodeScan mutator.
type CodeScanFunc func(context.Context, *ent.CodeScanMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CodeScanFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CodeScanMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CodeScanMutation", m)
	}
	return f(ctx, mv)
}

// The ComponentFunc type is an adapter to allow the use of ordinary
// function as Component mutator.
type ComponentFunc func(context.Context, *ent.ComponentMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ComponentFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ComponentMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ComponentMutation", m)
	}
	return f(ctx, mv)
}

// The ComponentUseFunc type is an adapter to allow the use of ordinary
// function as ComponentUse mutator.
type ComponentUseFunc func(context.Context, *ent.ComponentUseMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ComponentUseFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ComponentUseMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ComponentUseMutation", m)
	}
	return f(ctx, mv)
}

// The GitCommitFunc type is an adapter to allow the use of ordinary
// function as GitCommit mutator.
type GitCommitFunc func(context.Context, *ent.GitCommitMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f GitCommitFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.GitCommitMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.GitCommitMutation", m)
	}
	return f(ctx, mv)
}

// The LicenseFunc type is an adapter to allow the use of ordinary
// function as License mutator.
type LicenseFunc func(context.Context, *ent.LicenseMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LicenseFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.LicenseMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LicenseMutation", m)
	}
	return f(ctx, mv)
}

// The LicenseUseFunc type is an adapter to allow the use of ordinary
// function as LicenseUse mutator.
type LicenseUseFunc func(context.Context, *ent.LicenseUseMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LicenseUseFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.LicenseUseMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LicenseUseMutation", m)
	}
	return f(ctx, mv)
}

// The ProjectFunc type is an adapter to allow the use of ordinary
// function as Project mutator.
type ProjectFunc func(context.Context, *ent.ProjectMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProjectFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ProjectMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProjectMutation", m)
	}
	return f(ctx, mv)
}

// The ReleaseFunc type is an adapter to allow the use of ordinary
// function as Release mutator.
type ReleaseFunc func(context.Context, *ent.ReleaseMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ReleaseFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ReleaseMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ReleaseMutation", m)
	}
	return f(ctx, mv)
}

// The ReleaseEntryFunc type is an adapter to allow the use of ordinary
// function as ReleaseEntry mutator.
type ReleaseEntryFunc func(context.Context, *ent.ReleaseEntryMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ReleaseEntryFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ReleaseEntryMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ReleaseEntryMutation", m)
	}
	return f(ctx, mv)
}

// The RepoFunc type is an adapter to allow the use of ordinary
// function as Repo mutator.
type RepoFunc func(context.Context, *ent.RepoMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f RepoFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.RepoMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.RepoMutation", m)
	}
	return f(ctx, mv)
}

// The TestCaseFunc type is an adapter to allow the use of ordinary
// function as TestCase mutator.
type TestCaseFunc func(context.Context, *ent.TestCaseMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TestCaseFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.TestCaseMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TestCaseMutation", m)
	}
	return f(ctx, mv)
}

// The TestRunFunc type is an adapter to allow the use of ordinary
// function as TestRun mutator.
type TestRunFunc func(context.Context, *ent.TestRunMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TestRunFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.TestRunMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TestRunMutation", m)
	}
	return f(ctx, mv)
}

// The VulnerabilityFunc type is an adapter to allow the use of ordinary
// function as Vulnerability mutator.
type VulnerabilityFunc func(context.Context, *ent.VulnerabilityMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f VulnerabilityFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.VulnerabilityMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.VulnerabilityMutation", m)
	}
	return f(ctx, mv)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
//
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
//
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
//
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
//
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
