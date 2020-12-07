package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verifa/bubbly/env"
	"gopkg.in/h2non/gock.v1"
)

// TestQuery verifies that a call to c.Query is POSTed to the correct
// server route, and that the byte slice response matches the gock response
func TestQuery(t *testing.T) {
	tcs := []struct {
		desc         string
		query        string
		route        string
		responseCode int
		response     interface{}
		expectedCode int
	}{
		{
			desc: "basic query",
			query: `
			{
				  test_run {
					name
					test_set {
						name
						test_case {
							ID
							name
							status
							test_set_id
						}
					}
					repo_version_id
				}
			}
			`,
			route:        "/api/graphql",
			response:     `{"data":{"test_run":{"name":"run 1","repo_version_id":0,"test_set":[{"name":"set 1","test_case":[{"ID":1,"name":"case 1.1","status":"PASS","test_set_id":1},{"ID":2,"name":"case 1.2","status":"PASS","test_set_id":1},{"ID":3,"name":"case 1.3","status":"FAIL","test_set_id":1}]},{"name":"set 2","test_case":[{"ID":4,"name":"case 2.1","status":"FAIL","test_set_id":2},{"ID":5,"name":"case 2.2","status":"FAIL","test_set_id":2}]}]}}}`,
			responseCode: http.StatusOK,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()

			// Create a new server route for mocking a Bubbly server response
			gock.New(bCtx.ServerConfig.HostURL()).
				Post(tc.route).
				Reply(tc.responseCode).
				JSON(tc.response)

			c, err := New(bCtx)

			if err != nil {
				t.Errorf(err.Error())
			}

			byteRes, err := c.Query(bCtx, tc.query)

			t.Log(string(byteRes))
			assert.Equal(t, tc.response, string(byteRes))

		})
	}
}
