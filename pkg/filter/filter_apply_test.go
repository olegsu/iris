package filter_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/logger"
	"github.com/stretchr/testify/assert"

	mocks "github.com/olegsu/iris/pkg/filter/mocks"

	kubeMock "github.com/olegsu/iris/pkg/kube/mocks"
)

func Test_Filter_Apply(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name           string
		getFilter      func(string) filter.Filter
		args           args
		expectedResult bool
		expectedErr    bool
		filterType     string
	}{
		{
			filterType: filter.TypeAny,
			name:       "Should success when at least one inner filter is success",
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				mockedFactory := &mocks.Factory{}
				service := &mocks.Service{}
				mockedFilter1 := &mocks.Filter{}
				mockedFilter2 := &mocks.Filter{}

				anyFilter, _ := factory.Build(generateFilterJSON(t, "TEST", []string{"JSON-1", "JSON-2"}), service, nil)
				mockedFilter1.On("Apply", mock.Anything).Return(true, nil)
				mockedFilter2.On("Apply", mock.Anything).Return(false, nil)

				service.On("GetFilterByName", "JSON-1").Return(mockedFilter1, nil)
				service.On("GetFilterByName", "JSON-2").Return(mockedFilter2, nil)

				mockedFactory.On("Build", generateFilterJSON("jsonpath", "JSON-1", []string{"JSON-1"}), service, nil).Return(mockedFilter1, nil)
				mockedFactory.On("Build", generateFilterJSON("jsonpath", "JSON-2", []string{"JSON-2"}), service, nil).Return(mockedFilter2, nil)

				return anyFilter
			},
			args:           args{},
			expectedResult: true,
			expectedErr:    false,
		},
		{
			filterType: filter.TypeAny,
			name:       "Should fail when GetFilterByName returns an error",
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				anyFilter, _ := factory.Build(generateFilterJSON(t, "TEST", []string{"JSON-1", "JSON-2"}), service, nil)

				service.On("GetFilterByName", "JSON-1").Return(nil, errors.New("Error!"))

				return anyFilter
			},
			args:           args{},
			expectedResult: false,
			expectedErr:    true,
		},
		{
			filterType: filter.TypeAny,
			name:       "Should fail when filter.Apply returns an error",
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				mockedFactory := &mocks.Factory{}
				service := &mocks.Service{}
				mockedFilter := &mocks.Filter{}

				anyFilter, _ := factory.Build(generateFilterJSON(t, "TEST", []string{"JSON-1", "JSON-2"}), service, nil)
				mockedFilter.On("Apply", mock.Anything).Return(false, errors.New("Error!"))

				service.On("GetFilterByName", "JSON-1").Return(mockedFilter, nil)

				mockedFactory.On("Build", generateFilterJSON("jsonpath", "JSON-1", []string{"JSON-1"}), service, nil).Return(mockedFilter, nil)

				return anyFilter
			},
			args:           args{},
			expectedResult: false,
			expectedErr:    true,
		},
		{
			name:       "Should success mathcing json by exact value",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"value": "Value",
				}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: true,
			expectedErr:    false,
		},
		{
			name:       "Should failed mathcing json by exact value not match",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"value": "value",
				}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should fail mathcing json by exact value",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "not-mached",
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should success mathcing any value using regexp",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"regexp": ".", // match anything
				}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: true,
			expectedErr:    false,
		},
		{
			name:       "Should fail matching value using regexp",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"regexp": "VALUE", // match anything
				}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should fail when reading from json failed",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"regexp": "VALUE", // match VALUE
				}), service, nil)

				return jsonpathFilter
			},
			args:           args{},
			expectedResult: false,
			expectedErr:    true,
		},
		{
			name:       "Should fail when regexp is not valid",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"regexp": "(?:", // regex not valid
				}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: false,
			expectedErr:    true,
		},
		{
			name:       "Should fail when regexp or value not exist",
			filterType: filter.TypeJSONPath,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				jsonpathFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{}), service, nil)

				return jsonpathFilter
			},
			args: args{
				data: map[string]interface{}{
					"root": "Value",
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should success when namespace matched",
			filterType: filter.TypeNamespace,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				namespaceFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"namespace": "default",
				}), service, nil)

				return namespaceFilter
			},
			args: args{
				data: map[string]interface{}{
					"metadata": map[string]interface{}{
						"namespace": "default",
					},
				},
			},
			expectedResult: true,
			expectedErr:    false,
		},
		{
			name:       "Should fail when namespace not matched",
			filterType: filter.TypeNamespace,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				namespaceFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"namespace": "default",
				}), service, nil)

				return namespaceFilter
			},
			args: args{
				data: map[string]interface{}{
					"metadata": map[string]interface{}{
						"namespace": "iris",
					},
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should success when reason matched",
			filterType: filter.TypeReason,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				reasonFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"reason": "Scheduled",
				}), service, nil)

				return reasonFilter
			},
			args: args{
				data: map[string]interface{}{
					"reason": "Scheduled",
				},
			},
			expectedResult: true,
			expectedErr:    false,
		},
		{
			name:       "Should fail when reason not matched",
			filterType: filter.TypeReason,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}

				reasonFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"reason": "Scheduled",
				}), service, nil)

				return reasonFilter
			},
			args: args{
				data: map[string]interface{}{
					"reason": "Died",
				},
			},
			expectedResult: false,
			expectedErr:    false,
		},
		{
			name:       "Should success when reason matched",
			filterType: filter.TypeLabel,
			getFilter: func(t string) filter.Filter {
				// build the actual tested filter
				factory := filter.NewFactory(logger.New(nil))

				service := &mocks.Service{}
				kube := &kubeMock.Kube{}
				kube.On("ResourceByLabelsExist", mock.Anything, mock.Anything).Return(true, nil)
				reasonFilter, _ := factory.Build(generateFilterJSON(t, "TEST", map[string]interface{}{
					"app-version": "v1",
				}), service, kube)

				return reasonFilter
			},
			args: args{
				data: map[string]interface{}{
					"app-version": "v1",
				},
			},
			expectedResult: true,
			expectedErr:    false,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("Filter type:  %s  .%s", tt.filterType, tt.name)
		t.Run(name, func(t *testing.T) {
			fil := tt.getFilter(tt.filterType)
			got, err := fil.Apply(tt.args.data)
			assert.Equalf(t, tt.expectedErr, err != nil, "filter.Apply() error = %v, expectedErr %v", err, tt.expectedErr)

			if tt.expectedErr == false {
				assert.Equalf(t, tt.expectedResult, got, "filter.Apply() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
