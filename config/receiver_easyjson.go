// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA86519d3DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(in *jlexer.Lexer, out *ReceiverConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cloud_watch_region":
			out.CloudWatchRegion = string(in.String())
		case "cloud_watch_cluster_name":
			out.CloudWatchClusterName = string(in.String())
		case "aws_access_key_id":
			out.AWSAccessKeyID = string(in.String())
		case "aws_secret_access_key":
			out.AWSSecretAccessKey = string(in.String())
		case "http_request_timeout_millisecond":
			out.HTTPRequestTimeoutMillisecond = uint32(in.Uint32())
		case "metric_namespace":
			out.MetricNamespace = string(in.String())
		case "metric_name":
			out.MetricName = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA86519d3EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(out *jwriter.Writer, in ReceiverConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cloud_watch_region\":"
		out.RawString(prefix[1:])
		out.String(string(in.CloudWatchRegion))
	}
	{
		const prefix string = ",\"cloud_watch_cluster_name\":"
		out.RawString(prefix)
		out.String(string(in.CloudWatchClusterName))
	}
	{
		const prefix string = ",\"aws_access_key_id\":"
		out.RawString(prefix)
		out.String(string(in.AWSAccessKeyID))
	}
	{
		const prefix string = ",\"aws_secret_access_key\":"
		out.RawString(prefix)
		out.String(string(in.AWSSecretAccessKey))
	}
	{
		const prefix string = ",\"http_request_timeout_millisecond\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.HTTPRequestTimeoutMillisecond))
	}
	{
		const prefix string = ",\"metric_namespace\":"
		out.RawString(prefix)
		out.String(string(in.MetricNamespace))
	}
	{
		const prefix string = ",\"metric_name\":"
		out.RawString(prefix)
		out.String(string(in.MetricName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ReceiverConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA86519d3EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ReceiverConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA86519d3EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ReceiverConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA86519d3DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ReceiverConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA86519d3DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(l, v)
}
