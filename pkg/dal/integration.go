package dal

import (
	"encoding/json"
	"fmt"

	"github.com/olegsu/iris/pkg/destination"
	"github.com/olegsu/iris/pkg/filter"

	"k8s.io/api/core/v1"
)

type Integration struct {
	Name         string   `yaml:"name"`
	Filters      []string `yaml:"filters"`
	Destinations []string `yaml:"destinations"`
}

func (i *Integration) Exec(obj interface{}) (bool, error) {
	ev := obj.(*v1.Event)
	var j interface{}
	bytes, err := json.Marshal(&ev)
	if err != nil {
		return false, nil
	}
	json.Unmarshal(bytes, &j)
	result := true
	result = filter.IsFiltersMatched(GetDal().FilterService, i.Filters, j)
	if result == true {
		fmt.Printf("%s pass all checks, executing\n", i.Name)
		destination.Exec(GetDal().DestinationService, i.Destinations, obj)
	}
	return false, nil
}
