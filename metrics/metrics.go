package metrics

import (
	"fmt"
	"io"
	"net"
	"time"
)

// Metrics defines the functions supported for sending metrics.
type Metrics struct {
	Client Client
}

// New creates a new metrics client.
func New(client Client) *Metrics {
	return &Metrics{client}
}

// Timing logs timing information (in milliseconds).
func (m *Metrics) Timing(stat string, d time.Duration) error {
	_, err := fmt.Fprintf(m.Client, "%s:%d|ms\n", stat, d/time.Millisecond)
	return err
}

// Increment increments a specific stat counter by one.
func (m *Metrics) Increment(stat string) error {
	_, err := fmt.Fprintf(m.Client, "%s:%d|c\n", stat, 1)
	return err
}

// Decrement decrements a specific stat counter by one.
func (m *Metrics) Decrement(stat string) error {
	_, err := fmt.Fprintf(m.Client, "%s:%d|c\n", stat, -1)
	return err
}

// Client is an alias to io.Writer and is used with Metrics to separate the
// underlying connection.
type Client io.Writer

// The statsd type defines the relevant properties of a StatsD connection.
type statsd struct {
	addr string
	conn net.Conn
}

// Open establishes a udp connection to the StatsD server.
func (s *statsd) Open() error {
	conn, err := net.Dial("udp", s.addr)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

// Close closes the connection to the StatsD server.
func (s *statsd) Close() error {
	return s.conn.Close()
}

// Reload closes and re-opens the connection to the StatsD server. This is
// useful in long running processes which may need to reload their resources,
// typically following a HUP signal.
func (s *statsd) Reload() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return s.Open()
}

// Write satisfies the Client interface.
func (s *statsd) Write(b []byte) (int, error) {
	return s.conn.Write(b)
}

// Statsd creates a new client that can connect to a StatsD server. It's the
// users responsibility to Open the connecton following this call.
func Statsd(addr string) *statsd {
	return &statsd{addr: addr}
}
