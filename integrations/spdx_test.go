package integrations

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/store"
)

func TestSpdx(t *testing.T) {
	s, err := store.New(config.NewBubblyConfig())
	require.NoError(t, err)
	m, err := NewSPDXMonitor(WithStore(s))
	require.NoError(t, err)
	doErr := m.Do()
	require.NoError(t, doErr)
}
