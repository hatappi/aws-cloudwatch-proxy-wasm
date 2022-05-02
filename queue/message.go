//go:generate easyjson message.go

package queue

//easyjson:json
type Message struct {
	Host   string `json:"host"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Status string `json:"status"`
}
