package aws

import (
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/service"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/signature"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/httpcall"
)

type AWS struct {
	CloudWatch *service.CloudWatch
}

// New initializes AWS
func New(cwRegion, cwClusterName, accessKeyID, secretAccessKey string, timeoutMillisecond uint32) *AWS {
	httpClient := httpcall.New(timeoutMillisecond)
	signer := signature.New(accessKeyID, secretAccessKey)

	return &AWS{
		CloudWatch: service.NewCloudWatch(cwRegion, cwClusterName, httpClient, signer),
	}
}
