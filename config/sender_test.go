//go:build proxytest

package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadSenderConfig(t *testing.T) {
	testCases := map[string]struct {
		inputData []byte

		expectedConfig *SenderConfig
		wantErr        bool
	}{
		"valid inputData": {
			inputData: []byte(`{"receiver_vm_id": "foo"}`),
			expectedConfig: &SenderConfig{
				ReceiverVMID: "foo",
			},
			wantErr: false,
		},
		"nil inputData": {
			inputData:      nil,
			expectedConfig: nil,
			wantErr:        true,
		},
		"empty inputData": {
			inputData:      []byte(`{}`),
			expectedConfig: nil,
			wantErr:        true,
		},
		"invalid inputData": {
			inputData:      []byte(`foo`),
			expectedConfig: nil,
			wantErr:        true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c, err := LoadSenderConfig(tc.inputData)
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}

			if diff := cmp.Diff(c, tc.expectedConfig); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}
