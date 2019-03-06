package filter

import (
	"fmt"
	"strings"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/logger"
	"github.com/olegsu/iris/pkg/util"
)

// Factory interface builds filters
type Factory interface {
	Build(map[string]interface{}, Service, kube.Kube) (Filter, error)
}

type f struct {
	logger logger.Logger
}

// Build build actual filter and return Filter interface
func (_f *f) Build(json map[string]interface{}, s Service, k kube.Kube) (Filter, error) {
	if json["type"] != nil {
		t := strings.ToLower(json["type"].(string))
		var f Filter
		switch t {
		case TypeReason:
			f = &reasonFilter{
				baseFilter: baseFilter{
					logger: _f.logger,
				},
			}
			break
		case TypeNamespace:
			f = &namespaceFilter{
				baseFilter: baseFilter{
					logger: _f.logger,
				},
			}
			break
		case TypeJSONPath:
			f = &jsonPathFilter{
				baseFilter: baseFilter{
					logger: _f.logger,
				},
			}
			break
		case TypeLabel:
			f = &labelFilter{
				kube: k,
				baseFilter: baseFilter{
					logger: _f.logger,
				},
			}
			break
		case TypeAny:
			f = &anyFilter{
				Service: s,
				baseFilter: baseFilter{
					logger: _f.logger,
				},
			}
			break
		}
		if f == nil {
			return nil, fmt.Errorf("Type %s is not supported", json["type"])
		}
		util.MapToObjectOrDie(json, f, _f.logger)
		return f, nil
	} else {
		return nil, fmt.Errorf("Type passed to filter %v", json)
	}
}

// NewFactory create new factory
func NewFactory(logger logger.Logger) Factory {
	return &f{
		logger: logger,
	}
}
