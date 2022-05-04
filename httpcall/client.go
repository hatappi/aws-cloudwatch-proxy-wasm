package httpcall

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Callback represents callback function that has response result as arguments
type Callback func(headers http.Header, body []byte, err error)

// Client represents client for httpcall
type Client interface {
	Post(host, path string, header http.Header, body []byte, callback Callback) error
}

// New initializes client that meets Client interface
func New(timeoutMillisecond uint32) Client {
	return &client{
		timeoutMillisecond: timeoutMillisecond,
	}
}

type client struct {
	timeoutMillisecond uint32
}

// Post makes a POST request
func (c *client) Post(host, path string, header http.Header, body []byte, callback Callback) error {
	header.Set(":authority", host)
	header.Set(":method", "POST")
	header.Set(":path", path)
	header.Set("Accept", "application/json")

	if err := c.do(host, header, body, callback); err != nil {
		return fmt.Errorf("failed to make a request: %s", err)
	}

	return nil
}

func (c *client) do(host string, header http.Header, body []byte, callback Callback) error {
	cb := func(numHeaders, bodySize, numTrailers int) {
		hs, err := proxywasm.GetHttpCallResponseHeaders()
		if err != nil {
			callback(nil, nil, fmt.Errorf("failed to get response headers: %s", err))
			return
		}

		respHeader := make(http.Header)
		for _, v := range hs {
			respHeader.Add(strings.TrimLeft(v[0], ":"), v[1])
		}

		var body []byte
		if bodySize > 0 {
			body, err = proxywasm.GetHttpCallResponseBody(0, bodySize)
			if err != nil {
				callback(nil, nil, fmt.Errorf("failed to get response body: %s", err))
				return
			}
		}

		callback(respHeader, body, nil)
	}

	var reqHeaders [][2]string
	for k := range header {
		reqHeaders = append(reqHeaders, [2]string{k, header.Get(k)})
	}

	if _, err := proxywasm.DispatchHttpCall(host, reqHeaders, body, nil, c.timeoutMillisecond, cb); err != nil {
		return err
	}

	return nil
}
