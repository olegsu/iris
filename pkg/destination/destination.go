package destination

import (
	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
)

const (
	TypeDefault   = ""
	TypeCodefresh = "codefresh"
)

type Destination interface {
	GetName() string
	Exec(payload interface{})
	GetType() string
}

func NewDestination(json map[string]interface{}, k kube.Kube) Destination {
	var destination Destination
	if json["type"] != nil {
		destination = &codefreshDestination{}
	} else {
		destination = &defaultDestination{}
	}
	util.MapToObjectOrDie(json, destination)
	return destination
}

type baseDestination struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (d *baseDestination) GetName() string {
	return d.Name
}

func (d *baseDestination) GetType() string {
	return d.Type
}
