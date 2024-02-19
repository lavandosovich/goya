package tests

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	checkFunc func(resp *http.Response),
) (int, string) {

	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	checkFunc(resp)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
