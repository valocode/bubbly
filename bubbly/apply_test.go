package bubbly

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/config"
	"gopkg.in/h2non/gock.v1"
)

// TestApply simply validates that a given directory containing bubbly
// configuration including a pipeline_run will result in a POST of data to
// the bubbly server.
// See client/load_test.go for actual evaluation of the loading using
// the gofight package.
func TestApply(t *testing.T) {
	host := "localhost"
	port := "8080"
	auth := false
	token := ""
	hostURL := host + ":" + port

	// Create a new server route for mocking a Bubbly server response
	gock.New(hostURL).
		Post("/alpha1/upload").
		Reply(200).
		JSON(map[string]interface{}{"status": "uploaded"})

	sc := config.ServerConfig{
		Port:  port,
		Host:  "http://" + host,
		Auth:  auth,
		Token: token,
	}

	err := Apply("./testdata/sonarqube", sc)

	assert.NoError(t, err, "Failed to apply resource")

}
