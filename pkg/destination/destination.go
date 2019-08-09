package destination

import (
	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/logger"
	"github.com/olegsu/iris/pkg/util"
	"os"
)

const (
	TypeDefault   = ""
	TypeCodefresh = "codefresh"
)

var (
	UserAgent = "iris/" + os.Getenv("IRIS_VERSION")
)

type Destination interface {
	GetName() string
	Exec(payload interface{})
	GetType() string
}

func NewDestination(json map[string]interface{}, k kube.Kube, logger logger.Logger) Destination {
	var destination Destination
	if json["type"] != nil {
		destination = &codefreshDestination{
			baseDestination: baseDestination{
				logger: logger,
			},
		}
	} else {
		destination = &defaultDestination{
			baseDestination: baseDestination{
				logger: logger,
			},
		}
	}
	util.MapToObjectOrDie(json, destination, logger)
	return destination
}

type baseDestination struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	logger logger.Logger
}

func (d *baseDestination) GetName() string {
	return d.Name
}

func (d *baseDestination) GetType() string {
	return d.Type
}
