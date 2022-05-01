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

// NewHttpContext initializes senderHTTPContext
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	queueID, err := proxywasm.ResolveSharedQueue("receiver", constant.QueueName)
	if err != nil {
		proxywasm.LogCriticalf("failed to resolve queue %s: %v", constant.QueueName, err)
	}

	return &senderHTTPContext{
		contextID: contextID,
		queueID:   queueID,
	}
}

type senderHTTPContext struct {
	types.DefaultHttpContext

	contextID uint32
	queueID   uint32
}

// OnHttpStreamDone is called when the host environment is done processing
func (ctx *senderHTTPContext) OnHttpStreamDone() {
	authority, err := proxywasm.GetHttpRequestHeader(":authority")
	if err != nil {
		proxywasm.LogErrorf("failed to get authority header: %v", err)
		return
	}

	method, err := proxywasm.GetHttpRequestHeader(":method")
	if err != nil {
		proxywasm.LogErrorf("failed to get method header: %v", err)
		return
	}

	path, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		proxywasm.LogErrorf("failed to get path header: %v", err)
		return
	}

	status, err := proxywasm.GetHttpResponseHeader(":status")
	if err != nil {
		proxywasm.LogErrorf("failed to get status header: %v", err)
		return
	}

	proxywasm.LogDebugf("authority: %s, method: %s, path: %s, status: %s", authority, method, path, status)

	if err := proxywasm.EnqueueSharedQueue(ctx.queueID, []byte("test")); err != nil {
		proxywasm.LogCriticalf("failed to queue: %v", err)
	}
}
