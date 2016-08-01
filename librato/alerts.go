package librato

import (
	"errors"
	"fmt"
	"net/url"
)

type Alert struct {
	ID           int64                  `json:"id,omitempty"` // For responses. Do not include when posting.
	Name         string                 `json:"name"`
	Version      int64                  `json:"version"`
	Conditions   []Condition            `json:"conditions,omitempty"`
	Services     []Service              `json:"services,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Active       bool                   `json:"active,omitempty"`
	RearmSeconds int64                  `json:"rearm_seconds,omitempty"`
}

type Condition struct {
	Type            string                 `json:"type"`
	MetricName      string                 `json:"metric_name"`
	Source          string                 `json:"source,omitempty"`
	Threshold       float64                `json:"threshold,omitempty"`
	SummaryFunction string                 `json:"summary_function,omitempty"`
	Duration        int64                  `json:"duration"`
	DetectReset     bool                   `json:"detect_reset,omitempty"`
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
}

type AlertsResponse struct {
	Query  *QueryResponse `json:"query"`
	Alerts []*Alert       `json:"alerts"`
}

// PostAlert creates a new alert.
func (cli *Client) PostAlert(a *Alert) (int64, error) {
	if err := cli.request("POST", alertsURL, a, a); err != nil {
		return 0, err
	}
	return a.ID, nil
}

// GetAlerts returns a list of alerts
func (cli *Client) GetAlerts(name string, page *Pagination) (*AlertsResponse, error) {
	params := url.Values{}
	params.Set("version", "2")

	if name != "" {
		params.Set("name", name)
	}

	var a AlertsResponse
	return &a, cli.request("GET", alertsURL+"?"+page.toParams(params).Encode(), nil, &a)
}

//GetAlert returns an alert given its id
func (cli *Client) GetAlert(id int64) (*Alert, error) {
	var a Alert
	return &a, cli.request("GET", fmt.Sprintf("%s/%d", alertsURL, id), nil, &a)
}

//UpdateAlert updates an alert
func (cli *Client) PutAlert(a *Alert) error {
	if a.ID == 0 {
		return errors.New("Id cannot be 0")
	}
	return cli.request("PUT", fmt.Sprintf("%s/%d", alertsURL, a.ID), a, nil)
}

//DeleteAlert returns an alert given its id
func (cli *Client) DeleteAlert(id int64) error {
	return cli.request("DELETE", fmt.Sprintf("%s/%d", alertsURL, id), nil, nil)
}
