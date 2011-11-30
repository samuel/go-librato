package librato

import (
    "fmt"
    "http"
    "io"
    "os"
)

const(
    librato_metrics_api_url = "https://metrics-api.librato.com/v1/metrics.json"
    librato_metrics_users_api_url = "https://api.librato.com/v1/users.json"
    librato_metrics_services_api_url = "https://metrics-api.librato.com/v1/services.json"
)

func (q *QueryResponse) String() string {
    return fmt.Sprintf("{Found:%d Total:%d Offset:%d Length:%d}", q.Found, q.Total, q.Offset, q.Length)
}

func (met *Metrics) request(method string, url string, body io.Reader) (*http.Response, os.Error) {
 	req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    if method == "POST" {
        req.Header.Set("Content-Type", "application/json")
    }
    req.SetBasicAuth(met.Username, met.Token)
    res, err := http.DefaultClient.Do(req)
    return res, err
}

func (met *Metrics) get(url string) (*http.Response, os.Error) {
    return met.request("GET", url, nil)
}

func (met *Metrics) post(url string, body io.Reader) os.Error {
    res, err := met.request("POST", url, body)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return os.NewError(res.Status)
    }

    return nil
}
