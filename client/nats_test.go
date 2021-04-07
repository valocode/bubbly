package client

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/env"
)

// Just some random value that probably won't be in use. It can be changed
const TEST_PORT = 8131

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return natsserver.RunServer(&opts)
}

func TestNATS(t *testing.T) {
	bCtx := env.NewBubblyContext()
	bCtx.ClientConfig.ClientType = config.NATSClientType
	bCtx.ClientConfig.NATSAddr = fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)

	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	nc, err := nats.Connect(bCtx.ClientConfig.NATSAddr,
		nats.Name("Client Tests"),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			assert.NoError(t, err)
		}))
	require.NoErrorf(t, err, "nats connect")
	ec, err := nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)
	require.NoErrorf(t, err, "nats encoded connect")

	subjects := component.Subjects{component.StoreGetResourcesByKind,
		component.StorePostSchema, component.StoreQuery, component.StoreUpload}
	for _, sub := range subjects {
		ec.QueueSubscribe(string(sub), string(component.StoreQueue), func(m *nats.Msg) {
			reply, err := json.Marshal(component.Reply{})
			require.NoError(t, err)
			ec.Publish(m.Reply, reply)
		})
	}

	client, err := New(bCtx)
	require.NoErrorf(t, err, "failed to create NATS server")

	tables := core.Tables{
		core.Table{
			Name: "test_table",
		},
	}
	b, err := json.Marshal(tables)
	require.NoErrorf(t, err, "marshal tables")
	err = client.PostSchema(bCtx, b)
	require.NoErrorf(t, err, "post schema")

	// RESOURCE
	b, err = json.Marshal("test")
	require.NoError(t, err)
	err = client.PostResource(bCtx, b)
	require.NoError(t, err)
}
