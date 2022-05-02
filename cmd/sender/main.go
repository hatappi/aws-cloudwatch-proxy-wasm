package main

import (
	easyjson "github.com/mailru/easyjson"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

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

	config *config.SenderConfig
}

// OnPluginStart is called when the host environment starts the plugin
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("failed to get plugin config: %v", err)
		return types.OnPluginStartStatusFailed
	}

	config, err := config.LoadSenderConfig(configData)
	if err != nil {
		proxywasm.LogCriticalf("failed to load config: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.config = config

	return types.OnPluginStartStatusOK
}

// NewHttpContext initializes senderHTTPContext
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	queueID, err := proxywasm.ResolveSharedQueue(ctx.config.ReceiverVMID, constant.QueueName)
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

	message := queue.Message{
		Host:   authority,
		Method: method,
		Path:   path,
		Status: status,
	}

	mb, err := easyjson.Marshal(message)
	if err != nil {
		proxywasm.LogErrorf("failed to Marshal queue.Message: %v", err)
		return
	}

	cmb, err := proxywasm.CallForeignFunction("compress", mb)
	if err != nil {
		proxywasm.LogErrorf("failed to call compress function: %v", err)
		return
	}

	if err := proxywasm.EnqueueSharedQueue(ctx.queueID, cmb); err != nil {
		proxywasm.LogCriticalf("failed to queue: %v", err)
	}
}
