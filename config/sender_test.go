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

func TestMatchHosts_Contain(t *testing.T) {
	testCases := map[string]struct {
		matchHosts MatchHosts
		host       string
		expected   bool
	}{
		"contain": {
			matchHosts: []string{"example.com"},
			host:       "example.com",
			expected:   true,
		},
		"not contain": {
			matchHosts: []string{"example.com"},
			host:       "test.example.com",
			expected:   false,
		},
		"matchHosts is nil": {
			matchHosts: nil,
			host:       "test.example.com",
			expected:   false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			matched := tc.matchHosts.Contain(tc.host)
			if matched != tc.expected {
				t.Errorf("got: %v, expected: %v", matched, tc.expected)
			}
		})
	}
}
