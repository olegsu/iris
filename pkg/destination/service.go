package destination

import (
	"fmt"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
)

var d *dal

type dal struct {
	destinations []Destination
}

// Service
type Service interface {
	GetDestinationByName(string) (Destination, error)
}

// GetDestinationByName - finds a filters if exist
func (d *dal) GetDestinationByName(name string) (Destination, error) {
	var destination Destination
	if d.destinations == nil {
		return nil, fmt.Errorf("No destination %s", name)
	}
	for index := 0; index < len(d.destinations); index++ {
		filterName := d.destinations[index].GetName()
		if filterName == name {
			destination = d.destinations[index]
		}
	}
	if destination == nil {
		return nil, fmt.Errorf("%s destination not found", name)
	}
	return destination, nil
}

// NewService - creates net Dal from json array of filters
func NewService(destinationArray []map[string]interface{}, k kube.Kube) Service {
	tempDal := &dal{}
	for _, json := range destinationArray {
		f := NewDestination(json, k)
		tempDal.destinations = append(tempDal.destinations, f)
	}
	d = tempDal
	return tempDal
}

func Exec(serivce Service, names []string, payload interface{}) {
	for _, name := range names {
		dest, err := serivce.GetDestinationByName(name)
		if err != nil {
			util.EchoError(err)
		} else {
			dest.Exec(payload)
		}
	}
}
