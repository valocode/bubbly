package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceEvents(t *testing.T) {
	c := ResourceCreated

	require.Equal(t, "CREATED", c.String())
}
