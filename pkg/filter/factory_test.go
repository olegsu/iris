package filter

import (
	"reflect"
	"testing"

	kube "github.com/olegsu/iris/pkg/kube/mocks"
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		name string
		want Factory
	}{
		{
			name: "Get filter factory",
			want: &f{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_f_Build(t *testing.T) {

	tests := []struct {
		name         string
		buildJSONArg func() map[string]interface{}
		want         Filter
		wantErr      bool
	}{
		{
			name:    "Build filte with reason type",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type":   "reason",
					"name":   "filter-name",
					"reason": "reason",
				}
			},
			want: &reasonFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "reason",
				},
				Reason: "reason",
			},
		},
		{
			name:    "Build filte with namespace type",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type":      "namespace",
					"name":      "filter-name",
					"namespace": "default",
				}
			},
			want: &namespaceFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "namespace",
				},
				Namespace: "default",
			},
		},
		{
			name:    "Build filte with jsonpath type with value match",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type":  "jsonpath",
					"name":  "filter-name",
					"path":  "$.path",
					"value": "value",
				}
			},
			want: &jsonPathFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "jsonpath",
				},
				Path:  "$.path",
				Value: "value",
			},
		},
		{
			name:    "Build filte with jsonpath type with regexp match",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type":   "jsonpath",
					"name":   "filter-name",
					"path":   "$.path",
					"regexp": ".*",
				}
			},
			want: &jsonPathFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "jsonpath",
				},
				Path:   "$.path",
				Regexp: ".*",
			},
		},
		{
			name:    "Build filte with label type",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type": "labels",
					"name": "filter-name",
					"labels": map[string]string{
						"version": "v1",
					},
				}
			},
			want: &labelFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "labels",
				},
				Labels: map[string]string{
					"version": "v1",
				},
				kube: &kube.Kube{},
			},
		},
		{
			name:    "Build filte with any type",
			wantErr: false,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type": "any",
					"name": "filter-name",
					"filters": []string{
						"filter-1",
						"filter-2",
					},
				}
			},
			want: &anyFilter{
				baseFilter: baseFilter{
					Name: "filter-name",
					Type: "any",
				},
				Filters: []string{
					"filter-1",
					"filter-2",
				},
			},
		},
		{
			name:    "Return error when type not exit",
			wantErr: true,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{}
			},
			want: nil,
		},
		{
			name:    "Return error when type not supported",
			wantErr: true,
			buildJSONArg: func() map[string]interface{} {
				return map[string]interface{}{
					"type": "not-supported",
				}
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFactory().Build(tt.buildJSONArg(), &kube.Kube{})
			if (err != nil) != tt.wantErr {
				t.Errorf("f.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("f.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
