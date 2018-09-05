package filter_test

import (
	"github.com/olegsu/iris/pkg/filter"
)

func generateFilterJSON(t string, name string, data interface{}) map[string]interface{} {
	var r map[string]interface{}
	switch t {
	case filter.TypeAny:
		r = generateAnyFilter(name, data)
		break
	case filter.TypeJSONPath:
		r = generateJSONPathFilter(name, data)
		break
	case filter.TypeNamespace:
		r = generateNamespaceFilter(name, data)
		break
	case filter.TypeReason:
		r = generateReasonFilter(name, data)
		break
	case filter.TypeLabel:
		r = generateLabelFilter(name, data)
		break
	}
	return r
}

func generateAnyFilter(name string, data interface{}) map[string]interface{} {
	if name == "" {
		name = "generated"
	}
	return map[string]interface{}{
		"name":    name,
		"type":    filter.TypeAny,
		"filters": data,
	}
}

func generateJSONPathFilter(name string, data interface{}) map[string]interface{} {
	var res map[string]interface{}
	var path string
	if name == "" {
		name = "generated"
	}

	res = map[string]interface{}{
		"name": name,
		"type": filter.TypeJSONPath,
	}

	path = "$.root"

	casted, ok := data.(map[string]interface{})

	if ok {
		if casted["regexp"] != nil {
			res["regexp"] = casted["regexp"]
		} else if casted["value"] != nil {
			res["value"] = casted["value"]
		}
	}

	res["path"] = path
	return res
}

func generateNamespaceFilter(name string, data interface{}) map[string]interface{} {
	var res map[string]interface{}

	res = map[string]interface{}{
		"name": name,
		"type": filter.TypeNamespace,
	}

	casted, ok := data.(map[string]interface{})
	if ok {
		res["namespace"] = casted["namespace"]
	}
	return res
}

func generateReasonFilter(name string, data interface{}) map[string]interface{} {
	var res map[string]interface{}

	res = map[string]interface{}{
		"name": name,
		"type": filter.TypeReason,
	}

	casted, ok := data.(map[string]interface{})
	if ok {
		res["reason"] = casted["reason"]
	}
	return res
}

func generateLabelFilter(name string, data interface{}) map[string]interface{} {
	var res map[string]interface{}

	res = map[string]interface{}{
		"name":   name,
		"type":   filter.TypeLabel,
		"lables": data,
	}
	return res
}
