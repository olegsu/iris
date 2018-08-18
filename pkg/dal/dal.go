package dal

import (
	"github.com/olegsu/iris/pkg/destination"
	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
)

var dal *Dal

type Dal struct {
	DestinationService destination.Service
	Integrations       []*Integration
	FilterService      filter.Service
}

func GetDal() *Dal {
	return dal
}

func CreateDalFromBytes(bytes []byte, k kube.Kube) *Dal {
	d := &Dal{}
	var data map[string][]map[string]interface{}
	util.UnmarshalOrDie(bytes, &data)

	d.FilterService = filter.NewService(data["filters"], k)

	d.DestinationService = destination.NewService(data["destinations"], k)

	for _, integration := range data["integrations"] {
		i := &Integration{}
		util.MapToObjectOrDie(integration, i)
		d.Integrations = append(d.Integrations, i)
	}
	dal = d
	return d
}
