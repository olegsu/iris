package filter

import (
	"github.com/olegsu/iris/pkg/logger"
)

const (
	TypeReason    = "reason"
	TypeNamespace = "namespace"
	TypeJSONPath  = "jsonpath"
	TypeLabel     = "labels"
	TypeAny       = "any"
)

type Filter interface {
	GetName() string
	Apply(interface{}) (bool, error)
	GetType() string
}

func ApplyFilter(f Filter, obj interface{}) (bool, error) {
	return f.Apply(obj)
}

type baseFilter struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	logger logger.Logger
}

func (f *baseFilter) GetName() string {
	return f.Name
}

func (f *baseFilter) GetType() string {
	return f.Type
}
