package filter

type anyFilter struct {
	baseFilter `yaml:",inline"`
	Filters    []string `yaml:"filters"`
	Service    Service
}

func (f *anyFilter) Apply(data interface{}) (bool, error) {
	result := false
	var err error
	for _, name := range f.Filters {
		filter, err := f.Service.GetFilterByName(name)
		if err != nil {
			return false, err
		}
		res, err := filter.Apply(data)
		if err != nil {
			return false, err
		}
		if res == true {
			result = true
		}
	}
	return result, err
}
