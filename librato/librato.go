/*
Package librato is a library for Librato's metrics service API.

Example:

	metrics = &librato.Client{"login@email.com", "token"}
	metrics := &librato.Metrics{
		Counters: []librato.Metric{librato.Metric{"name", 123, "source"}},
		Gauges: []librato.Gauge{},
	}
	metrics.SendMetrics(metrics)
*/
package librato

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	userAgent = "go-librato/1.0"
)

const (
	annotationsURL = "https://metrics-api.librato.com/v1/annotations"
	metricsURL     = "https://metrics-api.librato.com/v1/metrics"
	servicesURL    = "https://metrics-api.librato.com/v1/services"
	usersURL       = "https://metrics-api.librato.com/v1/users"
)

type Sort string

const (
	Ascending  Sort = "asc"
	Descending Sort = "desc"
)

type QueryResponse struct {
	Found  int `json:"found"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type ErrTypes struct {
	Params  map[string]interface{} `json:"params"`
	Request []string               `json:"request"`
	System  []string               `json:"system"`
}

type ErrResponse struct {
	StatusCode int
	Errors     ErrTypes `json:"errors"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("librato: error %d: %+v", e.StatusCode, e.Errors)
}

type Client struct {
	Username string
	Token    string
}

func (cli *Client) request(method string, url string, req, res interface{}) error {
	if method == "GET" && req != nil {
		errors.New("librato: req must be nil for GET requests")
	}

	var body io.Reader
	if req != nil {
		buf := &bytes.Buffer{}
		body = buf
		if err := json.NewEncoder(buf).Encode(req); err != nil {
			return err
		}
	}

	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if httpReq != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.SetBasicAuth(cli.Username, cli.Token)
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode >= 400 {
		errRes := &ErrResponse{StatusCode: httpRes.StatusCode}
		if err := json.NewDecoder(httpRes.Body).Decode(errRes); err != nil {
			return err
		}
		return errRes
	}

	if res != nil {
		if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
			return err
		}
	}

	return nil
}

func pageParams(params url.Values, offset, length int, orderby string, sort Sort) url.Values {
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}
	if length > 0 {
		params.Set("length", strconv.Itoa(length))
	}
	if orderby != "" {
		params.Set("orderby", orderby)
	}
	if sort != "" {
		params.Set("sort", string(sort))
	}
	return params
}
