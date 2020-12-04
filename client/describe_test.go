// +build disabled

package client

import (
	"testing"

	"github.com/verifa/bubbly/env"

	"github.com/stretchr/testify/assert"

	"gopkg.in/h2non/gock.v1"
)

func TestDescribeResource(t *testing.T) {

	for _, c := range describeResourceCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			hostURL := c.sc.Host + ":" + c.sc.Port
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Get(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cl, err := NewClient(bCtx)

			if err != nil {
				t.Errorf(`Fail: %s: failed to create new test client: %s`, c.desc, err.Error())
			}
			actual, err := cl.DescribeResource(bCtx, c.rType, c.rName, c.rVersion)
			assert.Exactly(t, c.expected, actual)

		})
	}
}

func TestDescribeResourceGroup(t *testing.T) {

	for _, c := range describeResourceGroupCases {
		t.Run(c.desc, func(t *testing.T) {
			bCtx := env.NewBubblyContext()
			hostURL := c.sc.Host + ":" + c.sc.Port
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Get(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cl, err := NewClient(bCtx)

			if err != nil {
				t.Errorf(`Fail: %s: failed to create new test client: %s`, c.desc, err.Error())
			}
			actual, err := cl.DescribeResourceGroup(bCtx, c.rType, c.rVersion)
			assert.Exactly(t, c.expected, actual)

		})
	}
}
