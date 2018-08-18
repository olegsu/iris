package filter

import (
	"fmt"
	"strings"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
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

func NewFilter(json map[string]interface{}, k kube.Kube) (Filter, error) {
	if json["type"] != nil {
		t := strings.ToLower(json["type"].(string))
		var f Filter
		switch t {
		case TypeReason:
			f = &reasonFilter{}
			break
		case TypeNamespace:
			f = &namespaceFilter{}
			break
		case TypeJSONPath:
			f = &jsonPathFilter{}
			break
		case TypeLabel:
			f = &labelFilter{
				kube: k,
			}
			break
		case TypeAny:
			f = &anyFilter{}
			break
		}
		util.MapToObjectOrDie(json, f)
		return f, nil
	} else {
		return nil, fmt.Errorf("Type passed to filter %v\n", json)
	}
}

func ApplyFilter(f Filter, obj interface{}) (bool, error) {
	return f.Apply(obj)
}

type baseFilter struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (f *baseFilter) GetName() string {
	return f.Name
}

func (f *baseFilter) GetType() string {
	return f.Type
}
