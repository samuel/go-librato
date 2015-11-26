package librato

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func testId() string {
	// Create and seed the generator.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("test-alert-%d", r.Int63())
}

var alertsResponse = `
{
    "active": true,
    "attributes": {
        "runbook_url": "http://myco.com/runbooks/response_time"
    },
    "conditions": [
        {
            "metric_name": "web.nginx.response_time",
            "threshold": 200,
            "type": "above"
        }
    ],
    "id": 123,
    "name": "production.web.frontend.response_time",
    "services": [
        {
            "id": 849,
			"name": "campfire",
            "title": "Notify Campfire Room"
        }
    ],
    "version": 2
}
`

func TestUnmarshall(t *testing.T) {
	var a AlertsResponse

	err := json.Unmarshal([]byte(alertsResponse), &a)
	if err != nil {
		t.Fatalf("Error during unmarshal: %s", err.Error())
	}
}

func TestAlertCRUD(t *testing.T) {
	name := testId()
	cli := testClient(t)
	metricName := fmt.Sprintf("%s-metric", name)
	sourceName := fmt.Sprintf("%s-source", name)

	m := Metric{
		Name:        metricName,
		Value:       10,
		MeasureTime: time.Now().UTC().Unix(),
		Source:      sourceName}

	ms := Metrics{
		Gauges: []interface{}{m}}

	if err := cli.PostMetrics(&ms); err != nil {
		t.Fatal(err)
	}

	c := Condition{
		Type:       "above",
		Threshold:  200,
		Duration:   60,
		MetricName: metricName,
		Source:     sourceName}

	a := &Alert{
		Name:       name,
		Version:    2,
		Conditions: []Condition{c}}

	t.Logf("Posting alert %s", toJson(a))
	id, err := cli.PostAlert(a)
	if err != nil {
		t.Fatal(err)
	} else if id == 0 {
		t.Fatal("Id is 0")
	} else {
		t.Logf("Post response: %s", toJson(a))
	}

	ar, err := cli.GetAlerts(name, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(ar.Alerts) == 0 {
		t.Fatalf("Could not get alerts")
	}

	a2, err := cli.GetAlert(id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(a, a2) {
		t.Fatalf("expected %s but got %s", toJson(a), toJson(a2))
	}

	a.Conditions[0].Duration = 120

	if err := cli.PutAlert(a); err != nil {
		t.Fatal(err)
	}

	a2, err = cli.GetAlert(id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(a, a2) {
		t.Fatalf("expected %s but got %s", toJson(a), toJson(a2))
	}

	if err := cli.DeleteAlert(id); err != nil {
		t.Fatal(err)
	}

	a2, err = cli.GetAlert(id)
	if err == nil {
		t.Fatalf("Expected an error")
	}

	if err := cli.DeleteMetric(metricName); err != nil {
		t.Fatal(err)
	}
}
