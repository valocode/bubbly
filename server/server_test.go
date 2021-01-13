package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/verifa/bubbly/env"
)

func TestHealth(t *testing.T) {
	bCtx := env.NewBubblyContext()
	s := New(bCtx)

	router := s.setupRouter(bCtx)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

//
// func TestGetResource(t *testing.T) {
// 	bCtx := env.NewBubblyContext()
//
// 	s := New(bCtx)
//
// 	router := s.setupRouter(bCtx)
//
// 	w := httptest.NewRecorder()
// 	req, err := http.NewRequest(http.MethodGet,
// 		"/api/resource/default/pipelineRun/sonarqube", nil)
//
// 	require.NoError(t, err)
//
// 	router.ServeHTTP(w, req)
//
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
