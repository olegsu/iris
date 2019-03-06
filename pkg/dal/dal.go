package dal

import (
	"github.com/olegsu/iris/pkg/destination"
	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/logger"
	"github.com/olegsu/iris/pkg/util"
)

var dal *Dal

type Dal struct {
	DestinationService destination.Service
	Integrations       []*Integration
	FilterService      filter.Service
	Logger             logger.Logger
}

func GetDal() *Dal {
	return dal
}

func CreateDalFromBytes(bytes []byte, k kube.Kube, logger logger.Logger) *Dal {
	d := &Dal{
		Logger: logger,
	}
	var data map[string][]map[string]interface{}
	util.UnmarshalOrDie(bytes, &data, logger)

	d.FilterService = filter.NewService(filter.NewFactory(logger), data["filters"], k, logger)

	d.DestinationService = destination.NewService(data["destinations"], k, logger)

	for _, integration := range data["integrations"] {
		i := &Integration{}
		util.MapToObjectOrDie(integration, i, logger)
		i.logger = logger.New("Integration-Name", i.Name)
		d.Integrations = append(d.Integrations, i)
	}
	dal = d
	return d
}
