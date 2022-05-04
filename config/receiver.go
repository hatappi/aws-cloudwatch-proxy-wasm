//go:generate easyjson receiver.go
package config

import (
	"fmt"
	"os"

	easyjson "github.com/mailru/easyjson"
)

//easyjson:json
// ReceiverConfig represents config for receiver
type ReceiverConfig struct {
	CloudWatchRegion              string `json:"cloud_watch_region"`
	CloudWatchClusterName         string `json:"cloud_watch_cluster_name"`
	AWSAccessKeyID                string `json:"aws_access_key_id"`
	AWSSecretAccessKey            string `json:"aws_secret_access_key"`
	HTTPRequestTimeoutMillisecond uint32 `json:"http_request_timeout_millisecond"`
}

// LoadReceiverConfig parses data argument and returns ReceiverConfig
func LoadReceiverConfig(data []byte) (*ReceiverConfig, error) {
	config := &ReceiverConfig{
		AWSAccessKeyID:                os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:            os.Getenv("AWS_SECRET_ACCESS_KEY"),
		HTTPRequestTimeoutMillisecond: 5000,
	}

	if data != nil {
		if err := easyjson.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal receiver config: %s", err)
		}
	}

	// validation
	if config.CloudWatchRegion == "" {
		return nil, fmt.Errorf("cloud_watch_region is required")
	}
	if config.CloudWatchClusterName == "" {
		return nil, fmt.Errorf("cloud_watch_cluster_name is required")
	}
	if config.AWSAccessKeyID == "" {
		return nil, fmt.Errorf("aws_access_key_id is required")
	}
	if config.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("aws_secret_access_key is required")
	}

	return config, nil
}
