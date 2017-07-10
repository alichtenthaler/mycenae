package index

import "io"

// Metric refers to a timeseries metric
type Metric string

func (m Metric) String() string {
	return string(m)
}

// KVPair is a key-value pair for tags
type KVPair struct {
	Key   string
	Value string
}

// ID is a unique timeseries identifier
type ID uint64

// Backend defines the behaviour of the index
type Backend interface {
	// Add adds a new document to the index
	Add(Metric, []KVPair, ID) error
	// Query queries the underling index
	Query(Metric, []KVPair, []Filter) ResultSet

	// ListMetric lists all available metrics
	ListMetric(string) ([]string, error)
	// ListTagKeys lists all tag keys from the index
	ListTagKeys(string) ([]string, error)
	// ListTagValues lists all tag values given a tag key and metric, and a regexp
	ListTagValues(string, string) ([]string, error)

	// Store will save the content of the index in a file
	Store(io.Writer) error
	// Load will load data from a file
	Load(io.Reader) error
}
