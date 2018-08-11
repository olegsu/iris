package dal

import (
	"fmt"
	"regexp"

	"github.com/yalp/jsonpath"
)

type Filter struct {
	Name      string            `yaml:"name"`
	Type      string            `yaml:"type"`
	Namespace string            `yaml:"namespace"`
	Reason    string            `yaml:"reason"`
	Path      string            `yaml:"path"`
	Value     string            `yaml:"value"`
	Regexp    string            `yaml:"regexp"`
	Filters   []string          `yaml:"filters"`
	Labels    map[string]string `yaml:"labels"`
}

func (f *Filter) ConvertFilter() *Filter {
	if f.Reason != "" {
		return &Filter{
			Name:  f.Name,
			Type:  "jsonpath",
			Path:  "$.reason",
			Value: f.Reason,
		}
	}
	if f.Namespace != "" {
		return &Filter{
			Name:  f.Name,
			Type:  "jsonpath",
			Path:  "$.metadata.namespace",
			Value: f.Namespace,
		}
	}
	return f
}

func (f *Filter) Apply(data interface{}) (bool, error) {
	fmt.Printf("Applying filter %s\n", f.Name)
	f = f.ConvertFilter()
	if f.Type == "jsonpath" {
		res, err := applyJSONPathFilter(f, data)
		if err != nil {
			return false, err
		}
		return res, nil
	} else if f.Type == "any" {
		res, err := applyAnyFilter(f, data)
		if err != nil {
			return false, err
		}
		return res, nil
	} else if f.Type == "labels" {
		res, err := GetClientset("", true).FindResourceByLabels(data, f.Labels)
		if err != nil {
			return false, err
		}
		return res, nil
	}
	return false, nil
}

func applyAnyFilter(filter *Filter, data interface{}) (bool, error) {
	fmt.Printf("Any filter\n")
	var res = false
	for index := 0; index < len(filter.Filters); index++ {
		filter, err := dal.GetFilterByName(filter.Filters[index])
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

func applyJSONPathFilter(f *Filter, data interface{}) (bool, error) {
	path := f.Path
	actualValue, err := jsonpath.Read(data, path)
	if err != nil {
		fmt.Printf("Cant read json\n")
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
