package librato

import (
	"bytes"
	"json"
	"os"
)

type Metric struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Source string  `json:"source,omitempty"`
}

type MetricsFormat struct {
	Counters []Metric `json:"counters,omitempty"`
	Gauges   []Metric `json:"gauges,omitempty"`
}

type Metrics struct {
	Username string
	Token    string
}

func (met *Metrics) SendMetrics(metrics *MetricsFormat) os.Error {
	if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
		return nil
	}

	js, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return met.post(librato_metrics_api_url, bytes.NewBuffer(js))
}
