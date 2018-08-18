package filter

import (
	"fmt"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
)

var d *dal

type dal struct {
	filters []Filter
}

// Service
type Service interface {
	GetFilterByName(string) (Filter, error)
}

// GetFilterByName - finds a filters if exist
func (d *dal) GetFilterByName(name string) (Filter, error) {
	var f Filter
	if d.filters == nil {
		return nil, fmt.Errorf("No filters %s", name)
	}
	for index := 0; index < len(d.filters); index++ {
		filterName := d.filters[index].GetName()
		if filterName == name {
			f = d.filters[index]
		}
	}
	if f == nil {
		return nil, fmt.Errorf("%s filter not found", name)
	}
	return f, nil
}

// NewService - creates net Dal from json array of filters
func NewService(filterArray []map[string]interface{}, k kube.Kube) Service {
	tempDal := &dal{}
	for _, json := range filterArray {
		f, _ := NewFilter(json, k)
		tempDal.filters = append(tempDal.filters, f)
	}
	d = tempDal
	return tempDal
}

// IsFiltersMatched Go over all filters and apply each one on on data
// Return true is the data matched to all the filters
func IsFiltersMatched(service Service, requiredFilters []string, data interface{}) (bool, error) {
	matched := true
	var err error
	for _, f := range requiredFilters {
		var res bool
		filter, err := service.GetFilterByName(f)
		if err != nil {
			util.EchoError(err)
			matched = false
		}
		res, err = filter.Apply(data)
		if err != nil {
			util.EchoError(err)
			matched = false
		}
		if res == false {
			matched = false
		}
	}

	return matched, err
}
