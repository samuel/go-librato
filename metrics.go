/*
	Package librato is a library for Librato's metrics service API.

	Example:

		metrics = &librato.Metrics{"login@email.com", "token"}
		metrics := &librato.MetricsFormat{
			Counters: []librato.Metrics{librato.Metric{"name", 123, "source"}},
			Gauges: []librato.Metrics{},
		}
		metrics.SendMetrics(metrics)
*/
package librato

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	metricsApiUrl         = "https://metrics-api.librato.com/v1/metrics"
	metricsUsersApiUrl    = "https://api.librato.com/v1/users"
	metricsServicesApiUrl = "https://metrics-api.librato.com/v1/services"
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

// Crete and submit measurements for new or existing metrics.
// http://dev.librato.com/v1/post/metrics
func (met *Metrics) SendMetrics(metrics *MetricsFormat) error {
	if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
		return nil
	}

	js, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return met.post(metricsApiUrl, bytes.NewBuffer(js))
}

func (met *Metrics) request(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", userAgent)
	req.SetBasicAuth(met.Username, met.Token)
	res, err := http.DefaultClient.Do(req)
	return res, err
}

func (met *Metrics) get(url string) (*http.Response, error) {
	return met.request("GET", url, nil)
}

func (met *Metrics) post(url string, body io.Reader) error {
	res, err := met.request("POST", url, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
