package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestCallbackNoCode(t *testing.T) {
	w := httptest.NewRecorder()

	callback(w, httptest.NewRequest("", "/", nil), nil)

	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, w.Header().Get("Location"), "/")
}
