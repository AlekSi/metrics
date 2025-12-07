package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// NewCounter registers and returns new counter with the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
func NewCounter(name string) *Counter {
	return defaultSet.NewCounter(name)
}

// GetOrCreateCounter returns registered counter with the given name
// or creates new counter if the registry doesn't contain counter with
// the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// Performance tip: prefer NewCounter instead of GetOrCreateCounter.
func GetOrCreateCounter(name string) *Counter {
	return defaultSet.GetOrCreateCounter(name)
}

// Counter is a cumulative metric that represents a single monotonically increasing counter
// whose value can only increase or be reset to zero on restart.
//
// Counter is safe to use from concurrent goroutines.
// It must be passed by pointer; the value should not be copied.
type Counter struct {
	n    uint64
	help atomic.Pointer[string]
}

// Inc increments counter by 1.
func (c *Counter) Inc() {
	atomic.AddUint64(&c.n, 1)
}

// Dec decrements counter by 1.
//
// Using this method is discouraged, as counters should only increase.
// It is better to use [Gauge] instead.
func (c *Counter) Dec() {
	atomic.AddUint64(&c.n, ^uint64(0))
}

// Add increments counter by n.
func (c *Counter) Add(n int) {
	atomic.AddUint64(&c.n, uint64(n))
}

// AddInt64 increments counter by n.
func (c *Counter) AddInt64(n int64) {
	atomic.AddUint64(&c.n, uint64(n))
}

// Get returns the current counter value.
func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.n)
}

// Set sets counter value to n.
func (c *Counter) Set(n uint64) {
	atomic.StoreUint64(&c.n, n)
}

// SetHelp sets optional HELP metadata.
func (c *Counter) SetHelp(help string) {
	c.help.Store(&help)
}

func (c *Counter) Help() string {
	if h := c.help.Load(); h != nil {
		return *h
	}

	return ""
}

// metricType returns TYPE metadata for that metric: "counter".
func (c *Counter) metricType() string {
	return "counter"
}

// marshalTo marshals counter with the given prefix to w.
func (c *Counter) marshalTo(prefix string, w io.Writer) {
	v := c.Get()
	fmt.Fprintf(w, "%s %d\n", prefix, v)
}
