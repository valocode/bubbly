package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
)

type entryTypeMutation interface {
	Client() *ent.Client
	ReleaseID() (id int, exists bool)
	SetEntryID(id int)
	Type() string
}

func createReleaseEntry(ctx context.Context, m entryTypeMutation) error {
	var entryType releaseentry.Type
	switch m.Type() {
	case ent.TypeArtifact:
		entryType = releaseentry.TypeArtifact
	case ent.TypeCodeScan:
		entryType = releaseentry.TypeCodeScan
	case ent.TypeTestRun:
		entryType = releaseentry.TypeTestRun
	case ent.TypeCVEScan:
		entryType = releaseentry.TypeCveScan
	case ent.TypeLicenseScan:
		entryType = releaseentry.TypeLicenseScan
	default:
		return fmt.Errorf("unsupported release entry for type %s", m.Type())
	}

	client := m.Client()
	rID, ok := m.ReleaseID()
	if !ok {
		// Validator should already catch this as an artifact needs a release
		return errors.New("no release for artifact")
	}
	release, err := client.Release.Query().Where(release.ID(rID)).WithChecks().Only(ctx)
	if err != nil {
		// Validator should already catch this as an artifact needs a release
		return errors.New("no release for artifact")
	}
	// Create a release entry
	entry, err := client.ReleaseEntry.Create().
		SetRelease(release).
		SetType(entryType).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("error creating release entry for artifact: %w", err)
	}
	m.SetEntryID(entry.ID)
	return nil
}
