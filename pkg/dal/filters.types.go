package dal

import (
	"fmt"
	"regexp"

	"github.com/yalp/jsonpath"
)

var (
	TypeReason    = "reason"
	TypeNamespace = "namespace"
	TypeJSONPath  = "jsonpath"
	TypeLabel     = "labels"
	TypeAny       = "any"
)

type BaseFilter struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (f *BaseFilter) GetName() string {
	return f.Name
}

func (f *BaseFilter) GetType() string {
	return f.Type
}

type ReasonFilter struct {
	BaseFilter `yaml:",inline"`
	Reason     string `yaml:"reason"`
}

func (f *ReasonFilter) Apply(data interface{}) (interface{}, error) {
	jsonFilter := &JSONPathFilter{
		BaseFilter: BaseFilter{
			Name: f.GetName(),
			Type: f.GetType(),
		},
		Path:  "$.reason",
		Value: f.Reason,
	}
	return jsonFilter.Apply(data)
}

type NamespaceFilter struct {
	BaseFilter `yaml:",inline"`
	Namespace  string
}

func (f *NamespaceFilter) Apply(data interface{}) (interface{}, error) {
	jsonFilter := &JSONPathFilter{
		BaseFilter: BaseFilter{
			Name: f.GetName(),
			Type: f.GetType(),
		},
		Path:  "$.metadata.namespace",
		Value: f.Namespace,
	}
	return jsonFilter.Apply(data)
}

type JSONPathFilter struct {
	BaseFilter `yaml:",inline"`
	Path       string `yaml:"path"`
	Value      string `yaml:"value"`
	Regexp     string `yaml:"regexp"`
	Namespace  string `yaml:"namespace"`
}

func (f *JSONPathFilter) Apply(data interface{}) (interface{}, error) {
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

type LabelFilter struct {
	BaseFilter `yaml:",inline"`
	Labels     map[string]string `yaml:"labels"`
}

func (f *LabelFilter) Apply(data interface{}) (interface{}, error) {
	res, err := GetClientset("", true).FindResourceByLabels(data, f.Labels)
	if err != nil {
		return false, err
	}
	return res, nil
}

type AnyFilter struct {
	BaseFilter `yaml:",inline"`
	Filters    []string `yaml:"filters"`
}

func (f *AnyFilter) Apply(data interface{}) (interface{}, error) {
	fmt.Printf("Any filter\n")
	var res = false
	for index := 0; index < len(f.Filters); index++ {
		filter, err := dal.GetFilterByName(f.GetName())
		if err != nil {
			return false, err
		}
		res, err := filter.Apply(data)
		if err != nil {
			return false, err
		}
		if res == true {
			res = true
		}
	}
	return res, nil
}
