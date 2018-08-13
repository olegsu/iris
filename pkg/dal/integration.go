package dal

import (
	"encoding/json"
	"fmt"

	"github.com/olegsu/iris/pkg/util"
	"k8s.io/api/core/v1"
)

type Integration struct {
	Name         string   `yaml:"name"`
	Filters      []string `yaml:"filters"`
	Destinations []string `yaml:"destinations"`
}

func (i *Integration) Exec(obj interface{}) (bool, error) {
	fmt.Printf("Running integration %s\n", i.Name)
	ev := obj.(*v1.Event)
	var j interface{}
	bytes, err := json.Marshal(&ev)
	if err != nil {
		return false, nil
	}
	json.Unmarshal(bytes, &j)
	result := true
	for index := 0; index < len(i.Filters); index++ {
		filter := i.Filters[index]
		f, err := GetDal().GetFilterByName(filter)
		if err != nil {
			util.EchoError(err)
			return false, err
		}
		res, err := f.Apply(j)
		if err != nil {
			util.EchoError(err)
			return false, err
		}
		if res == false {
			result = false
			break
		}
	}
	if result == true {
		for index := 0; index < len(i.Destinations); index++ {
			dest := i.Destinations[index]
			destination, err := GetDal().GetDestinationByName(dest)
			if err != nil {
				util.EchoError(err)
			} else {
				destination.Exec(obj)
			}
		}
	}
	return false, nil
}
