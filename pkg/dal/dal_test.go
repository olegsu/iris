package dal

import (
	"testing"
)

func generateFakeDal(filters []Filter, destinations []Destination) *Dal {
	d := &Dal{}
	if filters != nil {
		d.Filters = filters
	}
	if destinations != nil {
		d.Destinations = destinations
	}
	return d
}

func generateFakeFilter(name string) Filter {
	return Filter{
		Name: name,
	}
}

func generateFakeDestination(name string) Destination {
	return Destination{
		Name: name,
	}
}

func TestNewDalFromFilePathWhenPathNoSet(t *testing.T) {
	d := NewDalFromFilePath("")
	if len(d.Destinations) > 0 || len(d.Filters) > 0 || len(d.Integrations) > 0 {
		t.Errorf("Should be empty")
	}
}

func TestNegativeGetFilterByName(t *testing.T) {
	var failTestTable = []struct {
		Name   string
		Dal    *Dal
		Input  string
		Output string
	}{
		{
			Name:   "Return not found error when filter not found",
			Dal:    generateFakeDal(nil, nil),
			Input:  "NotFound",
			Output: "NotFound filter not found",
		},
		{
			Name:   "Return not found error when dal not exist",
			Dal:    nil,
			Input:  "NotFound",
			Output: "NotFound filter not found",
		},
	}

	for _, test := range failTestTable {
		t.Logf("Running test: %s", test.Name)
		_, err := test.Dal.GetFilterByName(test.Input)
		if err == nil {
			t.Fail()
		} else if err.Error() != test.Output {
			t.Logf("Expected: %s\nActual: %s", err.Error(), test.Output)
			t.Fail()
		}
	}
	var successTestTable = []struct {
		Name   string
		Dal    *Dal
		Input  string
		Filter Filter
	}{
		{
			Name:   "Return the filter",
			Dal:    generateFakeDal([]Filter{generateFakeFilter("filter1")}, nil),
			Input:  "filter1",
			Filter: generateFakeFilter("filter1"),
		},
	}

	for _, test := range successTestTable {
		t.Logf("Running test: %s", test.Name)
		filter, err := test.Dal.GetFilterByName(test.Input)
		if err != nil {
			t.Logf("Error during get filter request, error: %s", err.Error())
			t.Fail()
		} else if filter.Name != test.Filter.Name {
			t.Logf("Expected: %s\nActual: %s", filter.Name, test.Filter.Name)
			t.Fail()
		}
	}
}

func TestGetDestinationByName(t *testing.T) {
	var failTestTable = []struct {
		Name   string
		Dal    *Dal
		Input  string
		Output string
	}{
		{
			Name:   "Return not found error when destination not found",
			Dal:    generateFakeDal(nil, nil),
			Input:  "NotFound",
			Output: "NotFound destination not found",
		},
		{
			Name:   "Return not found error when dal not exist",
			Dal:    nil,
			Input:  "NotFound",
			Output: "NotFound destination not found",
		},
	}

	for _, test := range failTestTable {
		t.Logf("Running test: %s", test.Name)
		_, err := test.Dal.GetDestinationByName(test.Input)
		if err == nil {
			t.Logf("Error should exist\n")
			t.Fail()
		} else if err.Error() != test.Output {
			t.Logf("Expected: %s\nActual: %s", err.Error(), test.Output)
			t.Fail()
		}
	}
	var successTestTable = []struct {
		Name        string
		Dal         *Dal
		Input       string
		Destination Destination
	}{
		{
			Name:        "Return the destination",
			Dal:         generateFakeDal(nil, []Destination{generateFakeDestination("destination1")}),
			Input:       "destination1",
			Destination: generateFakeDestination("destination1"),
		},
	}

	for _, test := range successTestTable {
		t.Logf("Running test: %s", test.Name)
		destination, err := test.Dal.GetDestinationByName(test.Input)
		if err != nil {
			t.Logf("Error during get destination request, error: %s", err.Error())
			t.Fail()
		} else if destination.Name != test.Destination.Name {
			t.Logf("Expected: %s\nActual: %s", destination.Name, test.Destination.Name)
			t.Fail()
		}
	}
}
