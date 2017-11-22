package influxLineProtocolOutput

import "testing"

func TestNew(t *testing.T) {
	slug := &MetricContainer{
		Name: "test_container",
	}

	c := New("test_container")

	if c == slug {
		t.Errorf("New did not return the correct values for a MetricContainer")
	}
}

func TestAddTags(t *testing.T) {
	slug := &MetricContainer{
		Name: "test_container",
		Tags: map[string]interface{}{
			"testTag1": 1,
			"TestTag2": "two",
		},
	}
	c := New("test_container")
	tag1 := make(map[string]interface{})
	tag1["testTag1"] = 1
	tag2 := map[string]interface{}{"TestTag2": "two"}

	c.AddTags(tag2)
	c.AddTags(tag1)

	if slug == c {
		t.Errorf(
			"Slug is not equal to c, got:\n%v\nwant:\n%v",
			slug,
			c,
		)
	}
}

func TestAddValues(t *testing.T) {
	slug := &MetricContainer{
		Name: "test_container",
		Values: map[string]interface{}{
			"testValue1": 1,
			"TestValue2": "two",
		},
	}
	c := New("test_container")
	tag1 := make(map[string]interface{})
	tag1["testValue1"] = 1
	tag2 := map[string]interface{}{"TestValue2": "two"}

	c.AddValues(tag2)
	c.AddValues(tag1)

	if slug == c {
		t.Errorf(
			"Slug is not equal to c, got:\n%v\nwant:\n%v",
			slug,
			c,
		)
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
		Tags: map[string]interface{}{
			"testTag1": 1,
			"TestTag2": "two",
		},
	}

	// Create the test container
	c := New("testing_container")
	c.AddTags(map[string]interface{}{
		"testTag1":   1,
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
