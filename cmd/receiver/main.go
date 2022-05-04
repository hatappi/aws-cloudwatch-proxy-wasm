package main

import (
	"time"

	easyjson "github.com/mailru/easyjson"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/hatappi/aws-cloudwatch-proxy-wasm/aws"
	aws_types "github.com/hatappi/aws-cloudwatch-proxy-wasm/aws/types"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/config"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/constant"
	"github.com/hatappi/aws-cloudwatch-proxy-wasm/queue"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// NewPluginContext initializes pluginxContext
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext

	config    *config.ReceiverConfig
	awsClient *aws.AWS
}

// OnPluginStart is called when the host environment starts the plugin
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	queueID, err := proxywasm.RegisterSharedQueue(constant.QueueName)
	if err != nil {
		proxywasm.LogCriticalf("failed to register queue %s: %v", constant.QueueName, err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("queue %s registered. queueID: %d", constant.QueueName, queueID)

	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("failed to get plugin config: %v", err)
		return types.OnPluginStartStatusFailed
	}

	config, err := config.LoadReceiverConfig(configData)
	if err != nil {
		proxywasm.LogCriticalf("failed to load config: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.config = config

	ctx.awsClient = aws.New(
		config.CloudWatchRegion,
		config.CloudWatchClusterName,
		config.AWSAccessKeyID,
		config.AWSSecretAccessKey,
		config.HTTPRequestTimeoutMillisecond,
	)

	return types.OnPluginStartStatusOK
}

// OnQueueReady is called when there is data available in the queue
func (ctx *pluginContext) OnQueueReady(queueID uint32) {
	data, err := proxywasm.DequeueSharedQueue(queueID)
	if err == types.ErrorStatusEmpty {
		proxywasm.LogDebugf("dequeued data from %d is empty", queueID)
		return
	}
	if err != nil {
		proxywasm.LogErrorf("failed to dequeue data from %d: %v", queueID, err)
		return
	}

	ud, err := proxywasm.CallForeignFunction("uncompress", data)
	if err != nil {
		proxywasm.LogErrorf("failed to call uncompress function: %v", err)
		return
	}

	var message queue.Message
	err = easyjson.Unmarshal(ud, &message)
	if err != nil {
		proxywasm.LogErrorf("failed to Unmarshal message: %v", err)
		return
	}

	err = ctx.awsClient.CloudWatch.PutMetricData(&aws_types.PutMetricDataInput{
		Namespace: ctx.config.MetricNamespace,
		MetricData: []aws_types.MetricDatum{
			{
				MetricName: ctx.config.MetricName,
				Timestamp:  time.Now().UTC().Unix(),
				Unit:       "Count",
				Value:      1,
				Dimensions: []aws_types.Dimension{
					{Name: "host", Value: message.Host},
					{Name: "method", Value: message.Method},
					{Name: "path", Value: message.Path},
					{Name: "status", Value: message.Status},
				},
			},
		},
	})
	if err != nil {
		proxywasm.LogErrorf("failed to call PutMetricData: %v", err)
		return
	}
}
