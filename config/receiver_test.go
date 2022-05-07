//go:build proxytest

package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadReceiverConfig(t *testing.T) {
	os.Unsetenv("ACPW_AWS_ACCESS_KEY_ID")
	os.Unsetenv("ACPW_AWS_SECRET_ACCESS_KEY")

	testCases := map[string]struct {
		inputData []byte

		expectedConfig *ReceiverConfig
		wantErr        bool
	}{
		"valid inputData": {
			inputData: []byte(`
				{
					"cloud_watch_region": "ap-northeast-1",
					"cloud_watch_cluster_name": "cloudwatch_cluster",
					"aws_access_key_id": "foo",
					"aws_secret_access_key": "bar",
					"http_request_timeout_millisecond": 100,
					"metric_namespace": "foo",
					"metric_name": "bar"
				}
			`),
			expectedConfig: &ReceiverConfig{
				CloudWatchRegion:              "ap-northeast-1",
				CloudWatchClusterName:         "cloudwatch_cluster",
				AWSAccessKeyID:                "foo",
				AWSSecretAccessKey:            "bar",
				HTTPRequestTimeoutMillisecond: 100,
				MetricNamespace:               "foo",
				MetricName:                    "bar",
			},
			wantErr: false,
		},
		"nil inputData": {
			inputData:      nil,
			expectedConfig: nil,
			wantErr:        true,
		},
		"cloud_watch_region is empty": {
			inputData: []byte(`
				{
					"cloud_watch_cluster_name": "cloudwatch_cluster",
					"aws_access_key_id": "foo",
					"aws_secret_access_key": "bar"
				}
			`),
			expectedConfig: nil,
			wantErr:        true,
		},
		"cloud_watch_cluster_name is empty": {
			inputData: []byte(`
				{
					"cloud_watch_region": "ap-northeast-1",
					"aws_access_key_id": "foo",
					"aws_secret_access_key": "bar"
				}
			`),
			expectedConfig: nil,
			wantErr:        true,
		},
		"aws_access_key_id is empty": {
			inputData: []byte(`
				{
					"cloud_watch_region": "ap-northeast-1",
					"cloud_watch_cluster_name": "cloudwatch_cluster",
					"aws_secret_access_key": "bar"
				}
			`),
			expectedConfig: nil,
			wantErr:        true,
		},
		"aws_secret_access_key is empty": {
			inputData: []byte(`
				{
					"cloud_watch_region": "ap-northeast-1",
					"cloud_watch_cluster_name": "cloudwatch_cluster",
					"aws_access_key_id": "foo"
				}
			`),
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
			c, err := LoadReceiverConfig(tc.inputData)
			if tc.wantErr != (err != nil) {
				t.Errorf("err: %v", err)
			}

			if diff := cmp.Diff(c, tc.expectedConfig); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}
