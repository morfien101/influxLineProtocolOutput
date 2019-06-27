package influxLineProtocolOutput

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	slug := &MetricContainer{
		Name: "foo",
	}

	c := New("foo")

	if reflect.TypeOf(slug) != reflect.TypeOf(c) {
		t.Errorf("New did not return the correct values for a MetricContainer")
	}
}

func TestAddTags(t *testing.T) {
	tags := map[string]string{
		"testTag1": "1",
		"TestTag2": "two",
	}
	c := New("test_container")
	tag1 := map[string]string{"testTag1": "1"}
	tag2 := map[string]string{"TestTag2": "two"}

	c.AddTags(tag2)
	c.AddTags(tag1)

	if err := c.ContainsTags(tags); err != nil {
		t.Errorf("AddTags did not add the correct values. Error: %s", err)
	}
}

func TestAddValues(t *testing.T) {
	values := map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
	}
	c := New("test_container")
	tag1 := map[string]interface{}{"testValue1": 1}
	tag2 := map[string]interface{}{"TestValue2": "two"}

	c.AddValues(tag2)
	c.AddValues(tag1)

	if err := c.ContainsValues(values); err != nil {
		t.Errorf("AddValues did not add the correct values. Error: %s", err)
	}
}

func TestAdd(t *testing.T) {
	values := map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
		"TestValue3": 1.2,
		"TestValue4": true,
	}
	tags := map[string]string{
		"testTag1": "1",
		"TestTag2": "two",
	}

	// Create the test container
	c := New("testing_container")
	c.Add(tags, values)

	if err := c.Contains(tags, values); err != nil {
		t.Errorf("Add did not add all the values. Error: %s", err)
	}
}

func TestOutput(t *testing.T) {
	// Create the test slug
	slug := &MetricContainer{
		Name: "test_container",
		Values: map[string]interface{}{
			"testValue1": 1,
			"TestValue2": "two",
			"TestValue3": 1.2,
			"TestValue4": true,
		},
		Tags: map[string]string{
			"testTag1": "1",
			"TestTag2": "two",
		},
	}

	// Create the test container
	c := New("testing_container")
	c.AddTags(map[string]string{
		"testTag1":   "1",
		"TestValue2": "two",
	})
	c.AddValues(map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
		"TestValue3": 1.2,
		"TestValue4": true,
	})

	// Now the test
	cOutput := c.Output()
	slugOutput := slug.Output()
	if cOutput == slugOutput {
		t.Errorf(
			"Slug output is not the same as c.\nWant:\n%s\nGot:\n%s",
			slugOutput,
			cOutput,
		)
	}
}

//	ContainsTags(map[string]string) error
//	ContainsValues(map[string]interface{}) error
//	Contains(map[string]string, map[string]interface{}) error
//	HasName(string) error
func TestTestingMethodsValid(t *testing.T) {
	// setup slugs to load into the tests
	slugValidName := "testing_container"
	slugValidValues := map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
		"TestValue3": 1.2,
		"TestValue4": true,
	}
	slugValidTags := map[string]string{
		"testTag1": "1",
		"TestTag2": "two",
	}
	slugTimestamp := time.Now().UnixNano()

	// Setup container to assert against.
	c := New("testing_container")
	c.AddTags(map[string]string{
		"testTag1": "1",
		"TestTag2": "two",
	})
	c.AddValues(map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
		"TestValue3": 1.2,
		"TestValue4": true,
	})

	c.SetTimeStamp(slugTimestamp)

	// Fire tests
	if err := c.HasName(slugValidName); err != nil {
		t.Errorf("HasName with a valid valid failed. Error: %s", err)
	}

	if err := c.HasTimestamp(slugTimestamp); err != nil {
		t.Errorf("HasTimestamp with a valid valid failed. Error: %s", err)
	}

	if err := c.ContainsTags(slugValidTags); err != nil {
		t.Errorf("ContainsTags with valid values failed. Error: %s", err)
	}

	if err := c.ContainsValues(slugValidValues); err != nil {
		t.Errorf("ContainsValues with valid values failed. Error: %s", err)
	}

	if err := c.Contains(slugValidTags, slugValidValues); err != nil {
		t.Errorf("Contains with valid values failed. Error: %s", err)
	}
}

func TestTestingMethodsInvalid(t *testing.T) {
	// setup slugs to load into the tests
	slugInvalidName := "testing_container"
	slugInvalidValues := map[string]interface{}{
		"testValue1": 10,
		"TestValue2": 2,
		"ImMissing":  1,
	}
	slugInvalidTags := map[string]string{
		"testTag1":  "1",
		"TestTag2":  "Three",
		"TestTag20": "two",
	}

	// Setup container to assert against.
	c := New("bean_counter")
	c.AddTags(map[string]string{
		"testTag1": "1",
		"TestTag2": "two",
	})
	c.AddValues(map[string]interface{}{
		"testValue1": 1,
		"TestValue2": "two",
		"TestValue3": 1.2,
		"TestValue4": true,
	})

	if err := c.HasName(slugInvalidName); err == nil {
		t.Error("HasName with invalid values failed to return an error.")
	}

	if err := c.ContainsTags(slugInvalidTags); err == nil {
		t.Error("ContainsTags with invalid values failed to return an error.")
	}

	if err := c.ContainsValues(slugInvalidValues); err == nil {
		t.Error("ContainsValues with invalid values failed to return an error.")
	}

	if err := c.Contains(slugInvalidTags, slugInvalidValues); err == nil {
		t.Error("Contains with invalid values failed to return an error.")
	}
}
