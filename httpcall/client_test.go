//go:build proxytest

package httpcall

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type testVmContext struct {
	types.DefaultVMContext
}

func (*testVmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &testPluginContext{}
}

type testPluginContext struct {
	types.DefaultPluginContext
}

func (*testPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &testHttpContext{}
}

type testHttpContext struct {
	types.DefaultHttpContext
}

func (ctx *testHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	client := New(1000)

	header := make(http.Header)
	header.Set("x-foo", "bar")

	err := client.Post("test_cluster", "/foo", header, nil, func(headers http.Header, body []byte, err error) {
		if err != nil {
			proxywasm.LogErrorf("failed to request: %v", err)
			proxywasm.ResumeHttpRequest()
			return
		}

		if err := proxywasm.SendHttpResponse(200, [][2]string{}, body, -1); err != nil {
			proxywasm.LogErrorf("failed to send local response: %v", err)
			proxywasm.ResumeHttpRequest()
			return
		}

		proxywasm.LogInfo("succeeded to request")

		proxywasm.ResumeHttpRequest()
	})
	if err != nil {
		proxywasm.LogErrorf("failed to make a request: %v", err)
	}

	return types.ActionPause
}

func TestPost(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&testVmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	contextID := host.InitializeHttpContext()

	// Call OnHttpRequestHeaders
	action := host.CallOnRequestHeaders(contextID, nil, false)
	if action != types.ActionPause {
		t.Errorf("action is not ActionPause: %d", action)
	}

	// Verify DispatchHttpCall is called
	actualAttrs := host.GetCalloutAttributesFromContext(contextID)
	expectedAttrs := []proxytest.HttpCalloutAttribute{
		{
			Upstream: "test_cluster",
			Headers: [][2]string{
				{"X-Foo", "bar"},
				{":authority", "test_cluster"},
				{":method", "POST"},
				{":path", "/foo"},
				{"Accept", "application/json"},
			},
			Trailers: [][2]string{},
			Body:     []byte{},
		},
	}
	opts := []cmp.Option{
		cmpopts.SortSlices(func(i, j [2]string) bool {
			return i[0] < j[0]
		}),
	}
	if diff := cmp.Diff(actualAttrs, expectedAttrs, opts...); diff != "" {
		t.Errorf("attributes of DispatchHttpCall mismatch (-got, +want)\n%s", diff)
	}

	// Call OnHttpCallResponse
	body := []byte("OK")
	headers := [][2]string{
		{"HTTP/1.1", "200 OK"},
		{"Date:", "Thu, 3 May 2022 00:00:00 GMT"},
		{"Content-Type", "application/json"},
		{"Content-Length", "2"},
	}
	host.CallOnHttpCallResponse(actualAttrs[0].CalloutID, headers, nil, body)

	// Verify local response
	actualLocalResponse := host.GetSentLocalResponse(contextID)
	expectedLocalResponse := &proxytest.LocalHttpResponse{
		StatusCode: 200,
		Data:       []byte("OK"),
		Headers:    [][2]string{},
		GRPCStatus: -1,
	}
	if diff := cmp.Diff(actualLocalResponse, expectedLocalResponse); diff != "" {
		t.Errorf("local response mismatch (-got, +want)\n%s", diff)
	}

	// Verify Envoy logs
	actualLogs := host.GetInfoLogs()
	expectedLogs := []string{"succeeded to request"}
	if diff := cmp.Diff(actualLogs, expectedLogs); diff != "" {
		t.Errorf("info logs mismatch (-got, +want)\n%s", diff)
	}
}
