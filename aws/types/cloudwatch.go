//go:generate easyjson cloudwatch.go

package types

//easyjson:json
type PutMetricDataInput struct {
	Namespace  string
	MetricData []MetricDatum
}

type MetricDatum struct {
	MetricName string
	Timestamp  int64
	Unit       string
	Value      float64
	Dimensions []Dimension
}

type Dimension struct {
	Name  string
	Value string
}
