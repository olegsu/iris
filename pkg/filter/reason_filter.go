package filter

type reasonFilter struct {
	baseFilter `yaml:",inline"`
	Reason     string `yaml:"reason"`
}

func (f *reasonFilter) Apply(data interface{}) (bool, error) {
	jsonFilter := &jsonPathFilter{
		baseFilter: baseFilter{
			Name: f.GetName(),
			Type: f.GetType(),
		},
		Path:  "$.reason",
		Value: f.Reason,
	}
	return jsonFilter.Apply(data)
}
