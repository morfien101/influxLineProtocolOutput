package influxLineProtocolOutput

import (
	"fmt"
	"strings"
	"sync"
)

// VERSION holds the version stamp
var VERSION = "0.1.0"

// MetricPrinter is used to output the metrics that have been gathered in a LineProtocol
type MetricPrinter interface {
	Output() string
	PrintOutput()
}

// MetricGather is used to save tags and values for output later
type MetricGather interface {
	AddTags(map[string]interface{})
	AddValues(map[string]string)
}

// Metric is used contains all the functions that metric can handle
type Metric interface {
	MetricPrinter
	MetricGather
}

// MetricContainer us used to hold the data
type MetricContainer struct {
	Name   string
	Tags   map[string]string
	Values map[string]interface{}
	sync.RWMutex
}

// New returns a fresh MetricsContainer pointer. It requires a string that is
// used as the name for the metric.
func New(name string) *MetricContainer {
	return &MetricContainer{
		Name:   name,
		Tags:   make(map[string]string),
		Values: make(map[string]interface{}),
	}
}

// AddTags will consume a map[string]interface{} and add it to the list of metrics
// that will be output later.
func (metric *MetricContainer) AddTags(data map[string]string) {
	metric.Lock()
	defer metric.Unlock()
	for key, value := range data {
		metric.Tags[key] = value
	}
}

// AddValues will consume a map[string]interface{} and add it to the list of metrics
// that will be output later.
func (metric *MetricContainer) AddValues(data map[string]interface{}) {
	metric.Lock()
	defer metric.Unlock()
	for key, value := range data {
		metric.Values[key] = value
	}
}

// Output will return the Line Protocol version of the metric
func (metric *MetricContainer) Output() string {
	// Lock for reading
	metric.RLock()
	defer metric.RLocker()

	// Name,Tags Values
	outformat := "%s,%s %s"

	// Create tags line
	var tagLine []string
	for k, v := range metric.Tags {
		tagLine = append(tagLine, fmt.Sprintf("%s=%v", k, v))
	}
	//fmt.Println("tagline:", tagLine)

	// Create values line
	var valueLine []string
	for k, v := range metric.Values {
		valueLine = append(valueLine, fmt.Sprintf("%s=%v", k, v))
	}
	//fmt.Println("valueLine:", valueLine)

	return fmt.Sprintf(
		outformat,
		metric.Name,
		strings.Join(tagLine, ","),
		strings.Join(valueLine, ","),
	)
}

// PrintOutput will convert the metric to a Line Protocol and then send it to STDOUT
func (metric *MetricContainer) PrintOutput() {
	fmt.Println(metric.Output())
}
