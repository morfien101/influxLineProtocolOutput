package influxLineProtocolOutput

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// VERSION holds the version stamp
var VERSION = "0.2.0"

// MetricPrinter is used to output the metrics that have been gathered in a LineProtocol
type MetricPrinter interface {
	Output() string
	PrintOutput()
}

// MetricGather is used to save tags and values for output later
type MetricGather interface {
	Add(map[string]string, map[string]interface{})
	AddTags(map[string]string)
	AddValues(map[string]interface{})
}

// Metric is used contains all the functions that metric can handle
type Metric interface {
	MetricPrinter
	MetricGather
}

// MetricTester is used to do assertions on the metric containers.
type MetricTester interface {
	ContainsTags(map[string]string) error
	ContainsValues(map[string]interface{}) error
	Contains(map[string]string, map[string]interface{}) error
	HasName(string) error
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

// Add consumes a map[string]string as your tags and second argument map[string]interface{} as your feilds.
// This is a short cut method if you have both your tags and fields ready to go.
func (metric *MetricContainer) Add(tags map[string]string, data map[string]interface{}) {
	metric.AddTags(tags)
	metric.AddValues(data)
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

func (metric *MetricContainer) name() string {
	metric.RLock()
	defer metric.RUnlock()
	return metric.Name
}

func (metric *MetricContainer) tags() map[string]string {
	metric.RLock()
	defer metric.RUnlock()
	return metric.Tags
}

func (metric *MetricContainer) values() map[string]interface{} {
	metric.RLock()
	defer metric.RUnlock()
	return metric.Values
}

// PrintOutput will convert the metric to a Line Protocol and then send it to STDOUT
func (metric *MetricContainer) PrintOutput() {
	fmt.Println(metric.Output())
}

// ContainsTags will assert that all the tags passed in are in the metric container.
// Tags keys and values are tested. An error with a string of the errors is passed back.
func (metric *MetricContainer) ContainsTags(testTags map[string]string) error {
	var invalid []string
	tags := metric.tags()
	for k, v := range testTags {
		if _, ok := tags[k]; ok {
			if tags[k] != v {
				invalid = append(invalid, fmt.Sprintf("Tag %s does not match stored value. Want: '%s'. Got:'%s'.", k, v, tags[k]))
			}
		} else {
			invalid = append(invalid, fmt.Sprintf("Tag %s was not found", k))
		}
	}
	if len(invalid) > 0 {
		return errors.New(strings.Join(invalid, ". "))

	}

	return nil
}

// ContainsValues will that all the values passed in are in the metric container.
// Values keys and values are tested. An err with the string of the errors is passed back.
func (metric *MetricContainer) ContainsValues(testValues map[string]interface{}) error {
	var invalid []string
	values := metric.values()
	for k, v := range testValues {
		if _, ok := values[k]; ok {
			if reflect.TypeOf(v) != reflect.TypeOf(values[k]) {
				// The types are not the same.
				invalid = append(
					invalid,
					fmt.Sprintf(
						"The types of the values are not the same for %s, want: %s, got:%s",
						k,
						fmt.Sprint(reflect.TypeOf(v)),
						fmt.Sprint(reflect.TypeOf(values[k])),
					),
				)
			} else {
				if v != values[k] {
					// The values are not the same
					invalid = append(
						invalid,
						fmt.Sprintf(
							"The values for %s do not match. Want: '%v'. Got: '%v'.",
							k,
							v,
							values[k],
						),
					)
				}
			}
		} else {
			// value is not there
			invalid = append(
				invalid,
				fmt.Sprintf(
					"The value %s is missing from the metrics container",
					k,
				),
			)
		}
	}

	if len(invalid) > 0 {
		return errors.New(strings.Join(invalid, ". "))
	}
	return nil
}

// Contains will call ContainsTags and ContainsValues on the passed in values. This is a short cut method.
func (metric *MetricContainer) Contains(testTags map[string]string, testValues map[string]interface{}) error {
	var err error
	var errorList []string

	err = metric.ContainsTags(testTags)
	if err != nil {
		errorList = append(errorList, fmt.Sprint(err))
	}

	err = metric.ContainsValues(testValues)
	if err != nil {
		errorList = append(errorList, fmt.Sprint(err))
	}

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	}

	return nil
}

// HasName will check that the metric container has the name passed in. If not it will return an error stating that fact.
func (metric *MetricContainer) HasName(name string) error {
	if metric.name() == name {
		return nil
	}
	return fmt.Errorf("The metric name does not match. Want: %s, got: %s", name, metric.name())
}
