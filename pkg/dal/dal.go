package dal

import (
	"fmt"

	"github.com/olegsu/iris/pkg/util"
)

var dal *Dal

type Dal struct {
	Destinations []Destination
	Integrations []Integration
	Filters      []Filter
}

func GetDal() *Dal {
	return dal
}

func NewDalFromFilePath(path string) *Dal {
	d := &Dal{}
	if path == "" {
		dal = d
		return dal
	}
	file := util.GetUtil().ReadFileOrDie(path)
	util.GetUtil().UnmarshalOrDie(file, d)
	dal = d
	return dal
}

func (dal *Dal) GetFilterByName(name string) (*Filter, error) {
	var f *Filter
	if dal == nil {
		return nil, fmt.Errorf("%s filter not found", name)
	}
	for index := 0; index < len(dal.Filters); index++ {
		if dal.Filters[index].Name == name {
			f = &dal.Filters[index]
		}
	}
	if f == nil {
		return nil, fmt.Errorf("%s filter not found", name)
	}
	return f, nil
}

func (dal *Dal) GetDestinationByName(name string) (*Destination, error) {
	var d *Destination
	if dal == nil {
		return nil, fmt.Errorf("%s destination not found", name)
	}
	for index := 0; index < len(dal.Destinations); index++ {
		if dal.Destinations[index].Name == name {
			d = &dal.Destinations[index]
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s destination not found", name)
	}
	return d, nil
}
