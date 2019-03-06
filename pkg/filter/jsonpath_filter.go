package filter

import (
	"regexp"

	"github.com/olegsu/iris/pkg/logger"
	"github.com/yalp/jsonpath"
)

type jsonPathFilter struct {
	baseFilter `yaml:",inline"`
	Path       string `yaml:"path"`
	Value      string `yaml:"value"`
	Regexp     string `yaml:"regexp"`
}

func (f *jsonPathFilter) Apply(data interface{}) (bool, error) {
	path := f.Path
	actualValue, err := jsonpath.Read(data, path)
	if err != nil {
		return false, err
	}
	if f.Value != "" {
		res := applyMatchValueFilter(f.Value, actualValue.(string), f.logger)
		return res, nil
	} else if f.Regexp != "" {
		res, err := applyRegexpFilter(f.Regexp, actualValue.(string), f.logger)
		if err != nil {
			return false, err
		}
		return res, nil
	} else {
		return false, nil
	}
}

func applyRegexpFilter(pattern string, value string, logger logger.Logger) (bool, error) {
	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false, err
	}
	if match == false {
		logger.Debug("JSON path does not match to regex", "pattern", pattern, "value", value)
		return false, nil
	}
	logger.Debug("JSON path match to regex", "pattern", pattern, "value", value)
	return true, nil
}

func applyMatchValueFilter(requiredValue string, actualValue string, logger logger.Logger) bool {
	if actualValue != requiredValue {
		logger.Debug("JSON path does not match", "requiredValue", requiredValue, "actualValue", actualValue)
		return false
	}
	logger.Debug("JSON path match", "requiredValue", requiredValue, "actualValue", actualValue)
	return true
}
