package librato

import (
	"time"

	"testing"
)

func TestQueryComposite(t *testing.T) {
	cli := testClient(t)

	res, err := cli.QueryComposite(`s("test_event","*")`, 60, time.Now().Add(-time.Hour*3), time.Time{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
	for _, m := range res.Measurements {
		t.Logf("\t%+v", m)
	}
}
