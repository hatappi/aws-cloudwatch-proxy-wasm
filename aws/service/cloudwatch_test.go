//go:build proxytest

package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/signature"
	aws_types "github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/types"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/httpcall"
)

type testHTTPClient struct {
	MockPost func(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error
}

func (thc *testHTTPClient) Post(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error {
	return thc.MockPost(cluster_name, host, path, header, body, callback)
}

type testSigner struct {
	MockSetSignatureV4Header func(header http.Header, config signature.SignatureV4Config) error
}

func (ts *testSigner) SetSignatureV4Header(header http.Header, config signature.SignatureV4Config) error {
	return ts.MockSetSignatureV4Header(header, config)
}

type testVmContext struct {
	types.DefaultVMContext

	httpClient httpcall.Client
	signer     signature.Signer
}

func (vmc *testVmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &testPluginContext{
		httpClient: vmc.httpClient,
		signer:     vmc.signer,
	}
}

type testPluginContext struct {
	types.DefaultPluginContext

	httpClient httpcall.Client
	signer     signature.Signer
}

func (pc *testPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &testHttpContext{
		httpClient: pc.httpClient,
		signer:     pc.signer,
	}
}

type testHttpContext struct {
	types.DefaultHttpContext

	httpClient httpcall.Client
	signer     signature.Signer
}

func (ctx *testHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	cw := NewCloudWatch("ap-northeast-1", "cw_cluster", ctx.httpClient, ctx.signer)

	err := cw.PutMetricData(&aws_types.PutMetricDataInput{})
	if err != nil {
		proxywasm.LogErrorf("failed to call PutMetricData: %v", err)
	}

	return types.ActionContinue
}

func TestPutMetricData(t *testing.T) {
	testCases := map[string]struct {
		httpClient httpcall.Client
		signer     signature.Signer

		expectedErrorLogCnt int
	}{
		"valid": {
			httpClient: &testHTTPClient{
				MockPost: func(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error {
					respHeader := make(http.Header)
					respHeader.Add("status", "200")

					callback(respHeader, nil, nil)

					return nil
				},
			},
			signer: &testSigner{
				MockSetSignatureV4Header: func(header http.Header, config signature.SignatureV4Config) error {
					return nil
				},
			},
			expectedErrorLogCnt: 0,
		},
		"request failed": {
			httpClient: &testHTTPClient{
				MockPost: func(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error {
					return errors.New("error")
				},
			},
			signer: &testSigner{
				MockSetSignatureV4Header: func(header http.Header, config signature.SignatureV4Config) error {
					return nil
				},
			},
			expectedErrorLogCnt: 1,
		},
		"request returns error": {
			httpClient: &testHTTPClient{
				MockPost: func(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error {
					respHeader := make(http.Header)
					respHeader.Add("status", "400")

					callback(respHeader, []byte("bad request"), nil)

					return nil
				},
			},
			signer: &testSigner{
				MockSetSignatureV4Header: func(header http.Header, config signature.SignatureV4Config) error {
					return nil
				},
			},
			expectedErrorLogCnt: 1,
		},
		"SetSignatureV4Header returns error": {
			httpClient: &testHTTPClient{
				MockPost: func(cluster_name, host, path string, header http.Header, body []byte, callback httpcall.Callback) error {
					respHeader := make(http.Header)
					respHeader.Add("status", "200")

					callback(respHeader, nil, nil)

					return nil
				},
			},
			signer: &testSigner{
				MockSetSignatureV4Header: func(header http.Header, config signature.SignatureV4Config) error {
					return errors.New("error")
				},
			},
			expectedErrorLogCnt: 1,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithVMContext(&testVmContext{
				httpClient: tc.httpClient,
				signer:     tc.signer,
			})
			host, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			contextID := host.InitializeHttpContext()
			host.CallOnRequestHeaders(contextID, nil, false)

			actualLogs := host.GetErrorLogs()
			if tc.expectedErrorLogCnt != len(actualLogs) {
				t.Errorf("log count is different. expected: %d, actual: %d (%v)", tc.expectedErrorLogCnt, len(actualLogs), actualLogs)
			}
		})
	}
}
