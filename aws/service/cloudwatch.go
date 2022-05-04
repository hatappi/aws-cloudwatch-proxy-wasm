package service

import (
	"fmt"
	"net/http"

	easyjson "github.com/mailru/easyjson"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/signature"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/types"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/httpcall"
)

const (
	cloudWatchAPIHost = "monitoring.ap-northeast-1.amazonaws.com"
)

type CloudWatch struct {
	region         string
	apiClusterName string
	httpClient     httpcall.Client
	signer         signature.Signer
}

// New initializes CloudWatch
func NewCloudWatch(region, apiClusterName string, httpClient httpcall.Client, signer signature.Signer) *CloudWatch {
	return &CloudWatch{
		region:         region,
		apiClusterName: apiClusterName,
		httpClient:     httpClient,
		signer:         signer,
	}
}

// PutMetricData sends the input to CloudWatch metrics
func (cw *CloudWatch) PutMetricData(input *types.PutMetricDataInput) error {
	body, err := easyjson.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to Marshal PutMetricDataInput: %s", err)
	}

	header := make(http.Header)
	header.Set("x-amz-target", "GraniteServiceVersion20100801.PutMetricData")

	err = cw.signer.SetSignatureV4Header(header, signature.SignatureV4Config{
		Region:  cw.region,
		Service: "monitoring",
		Method:  "POST",
		Host:    cloudWatchAPIHost,
		Path:    "/",
		Body:    body,
	})
	if err != nil {
		return fmt.Errorf("failed to set SignatureV4 header: %s", err)
	}

	err = cw.httpClient.Post(cw.apiClusterName, cloudWatchAPIHost, "/", header, body, func(respHeader http.Header, respBody []byte, err error) {
		if err != nil {
			proxywasm.LogErrorf("failed to request to PutMetricData: %v", err)
			return
		}

		if s := respHeader.Get("status"); s != "200" {
			proxywasm.LogErrorf("failed to request to PutMetricData: status: %s, body: %s", s, respBody)
			return
		}

		proxywasm.LogDebug("succeeded to request to PutMetricData")
	})
	if err != nil {
		return err
	}

	return nil
}
