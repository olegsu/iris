package filter

import (
	"fmt"
	"regexp"

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
		res := applyMatchValueFilter(f.Value, actualValue.(string))
		return res, nil
	} else if f.Regexp != "" {
		res, err := applyRegexpFilter(f.Regexp, actualValue.(string))
		if err != nil {
			return false, err
		}
		return res, nil
	} else {
		return false, nil
	}
}

func applyRegexpFilter(pattern string, value string) (bool, error) {
	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false, err
	}
	if match == false {
		return false, nil
	}
	fmt.Printf("JSON path match regex %s == %s\n", pattern, value)
	return true, nil
}

func applyMatchValueFilter(requiredValue string, actualValue string) bool {
	if actualValue != requiredValue {
		return false
	}
	fmt.Printf("JSON path match %s == %s\n", requiredValue, actualValue)
	return true
}
