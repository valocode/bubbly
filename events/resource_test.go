package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceEvents(t *testing.T) {
	c := ResourceCreatedUpdated

	require.Equal(t, "Created/Updated", c.String())
}
