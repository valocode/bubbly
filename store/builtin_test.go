package store

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/test"
	"github.com/zclconf/go-cty/cty"
)

func TestTime(t *testing.T) {
	bCtx := env.NewBubblyContext()
	resource := test.RunPostgresDocker(bCtx, t)
	bCtx.StoreConfig.PostgresAddr = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))

	s, err := New(bCtx)
	require.NoErrorf(t, err, "failed to initialize store")

	event := builtin.Event{
		Status: events.ResourceCreatedUpdated.String(),
		Time:   time.Now(),
	}

	tests := []struct {
		desc string
		dbs  core.DataBlocks
	}{
		{
			desc: "time type",
			dbs:  builtin.ToDataBlocks(event),
		},
		{
			desc: "time string",
			dbs: core.DataBlocks{core.Data{
				TableName: "_event",
				Fields: &core.DataFields{Values: map[string]cty.Value{
					"time": cty.StringVal(time.Now().Format(time.RFC3339)),
				}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err = s.Save(DefaultTenantName, tt.dbs)
			require.NoError(t, err)

			eventQuery := "{ _event(last:1) { time } }"
			var events builtin.Event_Wrap
			result, err := s.Query(DefaultTenantName, eventQuery)
			require.NoError(t, err)
			require.Empty(t, result.Errors)
			b, err := json.Marshal(result.Data)
			require.NoError(t, err)
			err = json.Unmarshal(b, &events)
			require.NoError(t, err)
		})
	}
}
