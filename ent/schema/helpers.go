package schema

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/ent/release"
	"github.com/valocode/bubbly/ent/releaseentry"
)

type entryTypeMutation interface {
	Client() *ent.Client
	ReleaseID() (id int, exists bool)
	Time() (r time.Time, exists bool)
	SetEntryID(id int)
	Type() string
}

type clientMutation interface {
	Client() *ent.Client
	Tx() (*ent.Tx, error)
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
	default:
		return fmt.Errorf("unsupported release entry for type %s", m.Type())
	}

	client := m.Client()
	rID, ok := m.ReleaseID()
	if !ok {
		// Validator should already catch this as an artifact needs a release
		return errors.New("no release for artifact")
	}
	release, err := client.Release.Query().Where(release.ID(rID)).Only(ctx)
	if err != nil {
		// Validator should already catch this as an artifact needs a release
		return errors.New("no release for artifact")
	}
	entryTime, ok := m.Time()
	if !ok {
		entryTime = time.Now()
	}
	// Create a release entry
	entry, err := client.ReleaseEntry.Create().
		SetRelease(release).
		SetTime(entryTime).
		SetType(entryType).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("error creating release entry for artifact: %w", err)
	}
	m.SetEntryID(entry.ID)
	return nil
}

func clientOrTxClient(m clientMutation) *ent.Client {
	if tx, err := m.Tx(); err == nil {
		// Means there is no transaction
		return tx.Client()
	} else {
		return m.Client()
	}
}
