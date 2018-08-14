package dal

import (
	"fmt"
	"strings"

	"github.com/olegsu/iris/pkg/util"
)

var dal *Dal

type Dal struct {
	Destinations []IDestination
	Integrations []*Integration
	Filters      []Ifilter
}

func GetDal() *Dal {
	return dal
}

func NewDalFromFilePath(path string) *Dal {
	if dal != nil {
		return dal
	}
	d := &Dal{}
	if path == "" {
		dal = d
		return dal
	}
	file := util.GetUtil().ReadFileOrDie(path)

	dal = createDalFromBytes(file)
	return dal
}

func NewDalFromConfigMap(configmapName string, configmapNamespace string) *Dal {
	if dal != nil {
		return dal
	}
	d := &Dal{}
	res, err := GetConfigmapData(configmapName, configmapNamespace)
	if err != nil {
		dal = d
		return d
	}
	bytes := []byte(res)
	dal = createDalFromBytes(bytes)
	return dal
}

func (dal *Dal) GetFilterByName(name string) (Ifilter, error) {
	var f Ifilter
	if dal == nil {
		return nil, fmt.Errorf("No filters %s", name)
	}
	for index := 0; index < len(dal.Filters); index++ {
		filterName := dal.Filters[index].GetName()
		if filterName == name {
			f = dal.Filters[index]
		}
	}
	if f == nil {
		return nil, fmt.Errorf("%s filter not found", name)
	}
	return f, nil
}

func (dal *Dal) GetDestinationByName(name string) (IDestination, error) {
	var d IDestination
	if dal == nil {
		return nil, fmt.Errorf("%s destination not found", name)
	}
	for index := 0; index < len(dal.Destinations); index++ {
		if dal.Destinations[index].GetName() == name {
			d = dal.Destinations[index]
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s destination not found", name)
	}
	return d, nil
}

func createDalFromBytes(bytes []byte) *Dal {
	d := &Dal{}
	var data map[string][]map[string]interface{}
	util.GetUtil().UnmarshalOrDie(bytes, &data)

	for _, filter := range data["filters"] {
		addFilterToDal(filter, d)
	}

	for _, destination := range data["destinations"] {
		addDestinationToDal(destination, d)
	}

	for _, integration := range data["integrations"] {
		i := &Integration{}
		util.GetUtil().MapToObjectOrDie(integration, i)
		d.Integrations = append(d.Integrations, i)
	}

	return d
}

func addDestinationToDal(destination map[string]interface{}, d *Dal) error {
	var dest IDestination
	if destination["type"] != nil {
		t := strings.ToLower(destination["type"].(string))
		switch t {
		case TypeCodefresh:
			dest = &CodefreshDestination{}
			break
		}
	} else {
		dest = &DefaultDestination{}
	}
	util.GetUtil().MapToObjectOrDie(destination, dest)
	d.Destinations = append(d.Destinations, dest)
	return nil
}

func addFilterToDal(filter map[string]interface{}, d *Dal) error {
	if filter["type"] != nil {
		t := strings.ToLower(filter["type"].(string))
		var f Ifilter
		switch t {
		case TypeReason:
			f = &ReasonFilter{}
			break
		case TypeNamespace:
			f = &NamespaceFilter{}
			break
		case TypeJSONPath:
			f = &JSONPathFilter{}
			break
		case TypeLabel:
			f = &LabelFilter{}
			break
		case TypeAny:
			f = &AnyFilter{}
			break
		}
		util.GetUtil().MapToObjectOrDie(filter, f)
		d.Filters = append(d.Filters, f)
		return nil
	} else {
		return fmt.Errorf("Type passed to filter %v\n", filter)
	}
}
