package filter

type namespaceFilter struct {
	baseFilter `yaml:",inline"`
	Namespace  string
}

func (f *namespaceFilter) Apply(data interface{}) (bool, error) {
	jsonFilter := &jsonPathFilter{
		baseFilter: baseFilter{
			Name:   f.GetName(),
			Type:   f.GetType(),
			logger: f.logger,
		},
		Path:  "$.metadata.namespace",
		Value: f.Namespace,
	}
	return jsonFilter.Apply(data)
}
