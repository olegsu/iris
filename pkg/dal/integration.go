package dal

import (
	"encoding/json"

	"github.com/olegsu/iris/pkg/destination"
	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/logger"

	v1 "k8s.io/api/core/v1"
)

type Integration struct {
	Name         string   `yaml:"name"`
	Filters      []string `yaml:"filters"`
	Destinations []string `yaml:"destinations"`
	logger       logger.Logger
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
	result = filter.IsFiltersMatched(GetDal().FilterService, i.Filters, j, i.logger)
	if result == true {
		i.logger.Debug("All checks are passed, executing", "name", i.Name)
		destination.Exec(GetDal().DestinationService, i.Destinations, obj, i.logger)
	}
	return false, nil
}
