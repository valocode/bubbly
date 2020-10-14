package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	router := SetupRouter()

	// Define Json structure
	jsonString := `{
		"name": "test1",
		"table": [
			{
				"name": "sub",
				"fields": [
					{
						"name": "testField1",
						"value": "testValue1"
					},
					{
						"name": "testField2",
						"value": "testValue2"
					}
				]
			}
		]
	}`

	var data Data
	// Create a Data struct based on the Json structure
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		panic(err.Error)
	}

	jsonReq, _ := json.Marshal(data)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload/alpha1/", bytes.NewBuffer(jsonReq))
	router.ServeHTTP(r, req)
	fmt.Println(r.Body)

	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "{\"status\":\"uploaded\"}", r.Body.String())
}
