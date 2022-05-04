//go:build proxytest

package signature

import (
	"net/http"
	"testing"
)

func TestSetSignatureV4Header(t *testing.T) {
	testCases := map[string]struct {
		config SignatureV4Config
		header http.Header

		wantErr           bool
		expectedHeaderCnt int
	}{
		"valid": {
			config: SignatureV4Config{
				Region:  "ap-northeast-1",
				Service: "monitoring",
				Method:  "POST",
				Host:    "monitoring.ap-northeast-1.amazonaws.com",
				Path:    "/",
				Body:    nil,
			},
			header:            http.Header{"X-foo": []string{"bar"}},
			wantErr:           false,
			expectedHeaderCnt: 4,
		},
		"header is nil": {
			config: SignatureV4Config{
				Region:  "ap-northeast-1",
				Service: "monitoring",
				Method:  "POST",
				Host:    "monitoring.ap-northeast-1.amazonaws.com",
				Path:    "/",
				Body:    nil,
			},
			header:            nil,
			wantErr:           true,
			expectedHeaderCnt: 0,
		},
	}

	signer := New("", "")

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			header := tc.header

			err := signer.SetSignatureV4Header(header, tc.config)
			if tc.wantErr != (err != nil) {
				t.Errorf("err: %v", err)
			}

			if tc.expectedHeaderCnt != len(header) {
				t.Errorf("header count is different. expected: %d, actual: %d", tc.expectedHeaderCnt, len(header))
			}
		})
	}
}
