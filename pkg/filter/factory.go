package filter

import (
	"fmt"
	"strings"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util"
)

// Factory interface builds filters
type Factory interface {
	Build(map[string]interface{}, Service, kube.Kube) (Filter, error)
}

type f struct{}

// Build build actual filter and return Filter interface
func (_f *f) Build(json map[string]interface{}, s Service, k kube.Kube) (Filter, error) {
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
			f = &anyFilter{
				Service: s,
			}
			break
		}
		if f == nil {
			return nil, fmt.Errorf("Type %s is not supported", json["type"])
		}
		util.MapToObjectOrDie(json, f)
		return f, nil
	} else {
		return nil, fmt.Errorf("Type passed to filter %v\n", json)
	}
}

// NewFactory create new factory
func NewFactory() Factory {
	return &f{}
}
