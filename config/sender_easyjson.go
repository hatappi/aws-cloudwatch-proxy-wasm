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

func easyjson1d751f33DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(in *jlexer.Lexer, out *SenderConfig) {
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
		case "receiver_vm_id":
			out.ReceiverVMID = string(in.String())
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
func easyjson1d751f33EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(out *jwriter.Writer, in SenderConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"receiver_vm_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ReceiverVMID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SenderConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1d751f33EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SenderConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1d751f33EncodeGithubComHatappiAwsCloudwatchProxyWasmConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SenderConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1d751f33DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SenderConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1d751f33DecodeGithubComHatappiAwsCloudwatchProxyWasmConfig(l, v)
}