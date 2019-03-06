package filter_test

import (
	"fmt"
	"testing"

	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/filter/mocks"
	kube "github.com/olegsu/iris/pkg/kube/mocks"
	"github.com/olegsu/iris/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_f_Build(t *testing.T) {
	tests := []struct {
		name         string
		buildJSONArg func() map[string]interface{}
		want         filter.Filter
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.New(nil)
			factory := filter.NewFactory(l)
			got, err := factory.Build(tt.buildJSONArg(), &mocks.Service{}, &kube.Kube{})
			assert.Equalf(t, tt.wantErr, err != nil, "f.Build() error = %v, wantErr %v", err, tt.wantErr)

			if tt.wantErr == false {
				message := fmt.Sprintf("Returned object does not implement Filter interface")
				assert.Implements(t, (*filter.Filter)(nil), got, message)
			}
		})
	}
}
