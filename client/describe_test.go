package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/h2non/gock.v1"
)

func TestDescribeResource(t *testing.T) {

	for _, c := range describeResourceCases {
		t.Run(c.desc, func(t *testing.T) {
			hostURL := c.sc.Host + ":" + c.sc.Port
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Get(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cl, err := NewClient(c.sc)

			if err != nil {
				t.Errorf(`Fail: %s: failed to create new test client: %s`, c.desc, err.Error())
			}
			actual, err := cl.DescribeResource(c.rType, c.rName, c.rVersion)
			assert.Exactly(t, c.expected, actual)

		})
	}
}

func TestDescribeResourceGroup(t *testing.T) {

	for _, c := range describeResourceGroupCases {
		t.Run(c.desc, func(t *testing.T) {
			hostURL := c.sc.Host + ":" + c.sc.Port
			// Create a new server route for mocking a Bubbly server response
			gock.New(hostURL).
				Get(c.route).
				Reply(c.responseCode).
				JSON(c.response)

			cl, err := NewClient(c.sc)

			if err != nil {
				t.Errorf(`Fail: %s: failed to create new test client: %s`, c.desc, err.Error())
			}
			actual, err := cl.DescribeResourceGroup(c.rType, c.rVersion)
			assert.Exactly(t, c.expected, actual)

		})
	}
}
