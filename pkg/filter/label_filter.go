package filter

import "github.com/olegsu/iris/pkg/kube"

type labelFilter struct {
	baseFilter `yaml:",inline"`
	Labels     map[string]string `yaml:"labels"`
	kube       kube.Kube
}

func (f *labelFilter) Apply(data interface{}) (bool, error) {
	return f.kube.FindResourceByLabels(data, f.Labels)
}
