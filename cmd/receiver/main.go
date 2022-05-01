package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/hatappi/aws-cloudwatch-proxy-wasm/constant"
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
}

// OnPluginStart is called when the host environment starts the plugin
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	queueID, err := proxywasm.RegisterSharedQueue(constant.QueueName)
	if err != nil {
		proxywasm.LogCriticalf("failed to register queue %s: %v", constant.QueueName, err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("queue %s registered. queueID: %d", constant.QueueName, queueID)

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

	proxywasm.LogErrorf("dequeued data: %s", data)
}
