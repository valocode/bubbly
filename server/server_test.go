package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/valocode/bubbly/env"
)

func TestHealth(t *testing.T) {
	bCtx := env.NewBubblyContext()
	s := New(bCtx)

	router := s.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
