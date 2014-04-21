package websocket

import (
	"strings"
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestServeWrongMethod(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://example.com/ws", strings.NewReader("key=value"))
	Serve(w, req)
	assert.Equal(t, w.Code, 405, "Bad request should return 'Method not allowed (405)'")
}

func TestServeWrongOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/ws", nil)
	Serve(w, r)
	assert.Equal(t, w.Code, 403, "Request from different domain should return 'Origin not allowed (403)'")
}

func TestServeNotAWebsocketRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/ws", nil)
	r.Header.Set("Origin", "http://example.com")
	Serve(w, r)
	assert.Equal(t, w.Code, 400, "Non websocket Requests from should return '(400)'")
}

