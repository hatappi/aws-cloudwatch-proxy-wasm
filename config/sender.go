//go:generate easyjson sender.go
package config

import (
	"fmt"

	easyjson "github.com/mailru/easyjson"
)

//easyjson:json
// SenderConfig represents config for sender
type SenderConfig struct {
	ReceiverVMID string `json:"receiver_vm_id"`
}

// LoadSenderConfig parses data argument and returns SenderConfig
func LoadSenderConfig(data []byte) (*SenderConfig, error) {
	config := &SenderConfig{}

	if data != nil {
		if err := easyjson.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sender config: %s", err)
		}
	}

	// validation
	if config.ReceiverVMID == "" {
		return nil, fmt.Errorf("receiver_vm_id is required")
	}

	return config, nil
}
