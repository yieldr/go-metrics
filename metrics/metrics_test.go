package metrics

import (
	"bytes"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	var client bytes.Buffer
	metrics := New(&client)

	for stat, wire := range map[string]string{
		"foo":     "foo:1|c\n",
		"foo.bar": "foo.bar:1|c\n",
	} {
		err := metrics.Increment(stat)
		if err != nil {
			t.Error(err)
		}
		if wire != client.String() {
			t.Errorf("expected %q but have %q", client.String(), wire)
		}
		client.Reset()
	}

	for _, test := range []struct {
		stat, wire string
		duration   time.Duration
	}{
		{"foo", "foo:1000|ms\n", 1 * time.Second},
		{"foo.bar", "foo.bar:200|ms\n", 200 * time.Millisecond},
	} {
		err := metrics.Timing(test.stat, test.duration)
		if err != nil {
			t.Error(err)
		}
		if test.wire != client.String() {
			t.Errorf("expected %q but have %q", client.String(), test.wire)
		}
		client.Reset()
	}
}
