// +build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/dataaccess"
)

func Test_Buntdb(t *testing.T) {
	dal, err := dataaccess.New(dataaccess.Config{
		Provider: dataaccess.ProviderType("buntdb"),
	})
	assert.NoError(t, err)

	assert.NoError(t, dal.P.Save("key", "coolValue"))
	resourceValue, err := dal.P.Query("key")
	assert.NoError(t, err)
	assert.Contains(t, resourceValue, "coolValue")
}

func Test_Etcd(t *testing.T) {
	dal, err := dataaccess.New(dataaccess.Config{
		Provider: dataaccess.ProviderType("etcd"),
	})
	assert.NoError(t, err)
	assert.NoError(t, dal.P.Save("key7", "test"))
	resourceValue, err := dal.P.Query("key7")
	assert.NoError(t, err)
	assert.Contains(t, resourceValue, "test")
}
