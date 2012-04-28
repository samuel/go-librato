package librato

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	metricsApiUrl         = "https://metrics-api.librato.com/v1/metrics"
	metricsUsersApiUrl    = "https://api.librato.com/v1/users"
	metricsServicesApiUrl = "https://metrics-api.librato.com/v1/services"
	userAgent             = "go-librato/0.5"
)

func (q *QueryResponse) String() string {
	return fmt.Sprintf("{Found:%d Total:%d Offset:%d Length:%d}", q.Found, q.Total, q.Offset, q.Length)
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
