# influxLineProtocolOutput
This is a libary package that will take in metrics and output them in Influx Line Protocol Output. It is designed to be used with external executables that will be invoked by Telegraf's exec plugin. This allows for easy to use metric gathering and outputting in a protocol that Telegraf can readily consume. Using this you can make plugins that don't have to live in the source code of Telegraf.

- https://docs.influxdata.com/influxdb/v0.9/write_protocols/line/
- https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec

## How it works
The documentation for this project can be generated using the GoDoc tool. However we will give a brief run down here of how it works.

There are 2 sections to this package. Building functions that you use to build out your application and testing functions that are helpers for testing.
### Building
There are 5 functions that you need to know about here.
The process is as follows, create a MetricContainer struct, add metrics then output at the end.
```golang
// New returns a fresh MetricsContainer pointer. It requires a string that is used as the name for the metric.
func New(name string) *MetricContainer
```

```golang
// Add consumes a map[string]string as your tags and second argument map[string]interface{} as your feilds. This is a short cut method if you have both your tags and fields ready to go.
func (metric *MetricContainer) Add(tags map[string]string, data map[string]interface{})
```

```golang
// AddTags will consume a map[string]interface{} and add it to the list of metrics that will be output later.
func (metric *MetricContainer) AddTags(data map[string]string)
```

```golang
// AddValues will consume a map[string]interface{} and add it to the list of metrics that will be output later.
func (metric *MetricContainer) AddValues(data map[string]interface{})
```

```golang
// Output will return the Line Protocol version of the metric
func (metric *MetricContainer) Output() string
```

```golang
// PrintOutput will convert the metric to a Line Protocol and then send it to STDOUT
func (metric *MetricContainer) PrintOutput()
```

### Testing
Testing has functions that can be called on a MetricContainer to assert data.

```golang
// Contains will call ContainsTags and ContainsValues on the passed in values. This is a short cut method.
func (metric *MetricContainer) Contains(testTags map[string]string, testValues map[string]interface{}) error
```

```golang
// ContainsTags will assert that all the tags passed in are in the metric container. Tags keys and values are tested. An error with a string of the errors is passed back.
func (metric *MetricContainer) ContainsTags(testTags map[string]string) error
```

```golang
// ContainsValues will that all the values passed in are in the metric container. Values keys and values are tested. An err with the string of the errors is passed back.
func (metric *MetricContainer) ContainsValues(testValues map[string]interface{}) error
```

```golang
// HasName will check that the metric container has the name passed in. If not it will return an error stating that fact.
func (metric *MetricContainer) HasName(name string) error

```

### Interfaces
The following interfaces are available for use.
```golang
// MetricGather is used to save tags and values for output later
type MetricGather interface {
    Add(map[string]string, map[string]interface{})
    AddTags(map[string]string)
    AddValues(map[string]interface{})
}
```

```golang
// MetricPrinter is used to output the metrics that have been gathered in a LineProtocol
type MetricPrinter interface {
    Output() string
    PrintOutput()
}
```

```golang
// MetricTester is used to do assertions on the metric containers.
type MetricTester interface {
    ContainsTags(map[string]string) error
    ContainsValues(map[string]interface{}) error
    Contains(map[string]string, map[string]interface{}) error
    HasName(string) error
}
```
