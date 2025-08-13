package utils

import (
	"net/http/httptest"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/storage"
)

// TestServer provides a test HTTP server with storage integration
type TestServer struct {
	*httptest.Server
	Storage *storage.Storage
	Router  *aphttp.Router
}

// NewTestServer creates a new test server with the given storage instance
func NewTestServer(storage *storage.Storage) *TestServer {
	router := aphttp.NewRouter(storage)
	server := httptest.NewServer(router.GetMux())

	return &TestServer{
		Server:  server,
		Storage: storage,
		Router:  router,
	}
}

// Close closes the test server
func (ts *TestServer) Close() {
	if ts.Server != nil {
		ts.Server.Close()
	}
}
