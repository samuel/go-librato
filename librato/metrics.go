package librato

import (
	"net/url"
	"strconv"
	"time"
)

type Metric struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	MeasureTime int64   `json:"measure_time,omitempty"`
	Source      string  `json:"source,omitempty"`
}

type Gauge struct {
	Name        string  `json:"name"`
	MeasureTime int64   `json:"measure_time,omitempty"`
	Source      string  `json:"source,omitempty"`
	Count       uint64  `json:"count"`
	Sum         float64 `json:"sum"`
	Max         float64 `json:"max,omitempty"`
	Min         float64 `json:"min,omitempty"`
	SumSquares  float64 `json:"sum_squares,omitempty"`
}

type Metrics struct {
	MeasureTime int64         `json:"measure_time,omitempty"`
	Source      string        `json:"source,omitempty"`
	Counters    []Metric      `json:"counters,omitempty"`
	Gauges      []interface{} `json:"gauges,omitempty"` // Values can be either Metric or Gauge
}

type ComposeResult struct {
	Compose      string         `json:"compose"`
	Measurements []*Measurement `json:"measurements"`
	Resolution   int            `json:"resolution"`
	Query        struct {
		NextTime *int64 `json:"next_time"`
	} `json:"query"`
}

type Measurement struct {
	Series []Value     `json:"series"`
	Metric *MetricInfo `json:"metric"`
	Source struct {
		Name string `json:"name"`
	} `json:"source"`
	Query struct {
		Metric string `json:"metric"`
		Source string `json:"source"`
	} `json:"query"`
	Period *int `json:"period,omitempty"`
	// Timeshift *string `json:"timeshift,omitempty"` // TODO: don't know the type
}

type MetricInfo struct {
	Name        string                 `json:"name"`
	DisplayName *string                `json:"display_name,omitempty"`
	Type        string                 `json:"type"`
	Attributes  map[string]interface{} `json:"attributes"`
	Description *string                `json:"description,omitempty"`
	Period      *int                   `json:"period,omitempty"`
	// SourceLag   *string           `json:"source_lag"` // TODO: don't know the type
}

type Value struct {
	Value       float64 `json:"value"`
	MeasureTime int64   `json:"measure_time"`
}

// PostMetrics submits measurements for new or existing metrics.
// http://dev.librato.com/v1/post/metrics
func (cli *Client) PostMetrics(metrics *Metrics) error {
	if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
		return nil
	}
	return cli.request("POST", metricsURL, metrics, nil)
}

func (cli *Client) QueryComposite(compose string, resolution int, startTime, endTime time.Time, count int) (*ComposeResult, error) {
	v := url.Values{
		"compose":    []string{compose},
		"resolution": []string{strconv.Itoa(resolution)},
		"start_time": []string{strconv.FormatInt(startTime.Unix(), 10)},
	}
	if !endTime.IsZero() {
		v.Set("end_time", strconv.FormatInt(endTime.Unix(), 10))
	}
	if count > 0 {
		v.Set("count", strconv.Itoa(count))
	}
	var res ComposeResult
	if err := cli.request("GET", metricsURL+"?"+v.Encode(), nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
