package filter_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/olegsu/iris/pkg/filter"
	filterMock "github.com/olegsu/iris/pkg/filter/mocks"
	"github.com/olegsu/iris/pkg/kube"
	kubeMock "github.com/olegsu/iris/pkg/kube/mocks"
	"github.com/stretchr/testify/mock"
)

func generateFilterAsJSONArray(len int) []map[string]interface{} {
	res := []map[string]interface{}{}
	for i := 0; i < len; i++ {
		toAdd := map[string]interface{}{
			"type": "jsonpath",
		}
		res = append(res, toAdd)
	}
	return res
}

func TestNewService(t *testing.T) {
	type args struct {
		factory     *filterMock.Factory
		filterArray []map[string]interface{}
		kube        kube.Kube
	}

	tests := []struct {
		name      string
		args      args
		callCount int
	}{
		{
			name: "Create service with no filters",
			args: args{
				factory:     &filterMock.Factory{},
				kube:        &kubeMock.Kube{},
				filterArray: generateFilterAsJSONArray(0),
			},
			callCount: 0,
		},
		{
			name: "Create service one filter",
			args: args{
				factory:     &filterMock.Factory{},
				kube:        &kubeMock.Kube{},
				filterArray: generateFilterAsJSONArray(1),
			},
			callCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.filterArray) > 0 {
				tt.args.factory.On("Build", tt.args.filterArray[0], tt.args.kube).Return(nil, nil)
			}
			filter.NewService(tt.args.factory, tt.args.filterArray, tt.args.kube)
			tt.args.factory.AssertNumberOfCalls(t, "Build", tt.callCount)
		})
	}
}

func TestIsFiltersMatched(t *testing.T) {
	type args struct {
		getServiceFn    func() filter.Service
		requiredFilters []string
		data            interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Return true when no filters passed",
			args: args{
				getServiceFn: func() filter.Service {
					return nil
				},
				requiredFilters: []string{},
				data:            nil,
			},
			want: true,
		},
		{
			name: "Return true when all filters return truly value",
			args: args{
				getServiceFn: func() filter.Service {
					s := &filterMock.Service{}
					f := &filterMock.Filter{}
					f.
						On("Apply", mock.Anything).
						Return(true, nil)
					s.
						On("GetFilterByName", mock.Anything).
						Return(f, nil)
					return s
				},
				requiredFilters: []string{"filter-name"},
				data:            nil,
			},
			want: true,
		},
		{
			name: "Return false when at least one filter return falsy value",
			args: args{
				getServiceFn: func() filter.Service {
					s := &filterMock.Service{}
					f := &filterMock.Filter{}
					f.
						On("Apply", mock.Anything).
						Return(false, nil)

					s.
						On("GetFilterByName", mock.Anything).
						Return(f, nil)
					return s
				},
				requiredFilters: []string{
					"filter-success",
					"filter-failed",
				},
				data: nil,
			},
			want: false,
		},
		{
			name: "Return false when filter returns an error",
			args: args{
				getServiceFn: func() filter.Service {
					s := &filterMock.Service{}
					f := &filterMock.Filter{}
					f.
						On("Apply", mock.Anything).
						Return(false, errors.New("Error!"))

					s.
						On("GetFilterByName", mock.Anything).
						Return(f, nil)
					return s
				},
				requiredFilters: []string{
					"filter-success",
					"filter-failed",
				},
				data: nil,
			},
			want: false,
		},
		{
			name: "Return false getFilterByName returns an error",
			args: args{
				getServiceFn: func() filter.Service {
					s := &filterMock.Service{}

					s.
						On("GetFilterByName", mock.Anything).
						Return(nil, errors.New("Error!"))
					return s
				},
				requiredFilters: []string{
					"filter-failed",
				},
				data: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.IsFiltersMatched(tt.args.getServiceFn(), tt.args.requiredFilters, tt.args.data)
			if got != tt.want {
				t.Errorf("IsFiltersMatched() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dal_GetFilterByName(t *testing.T) {
	type fields struct {
		filters []filter.Filter
	}
	type args struct {
		name string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               func(string) filter.Filter
		wantErr            bool
		getFilterJSONArray func() []map[string]interface{}
		getFactory         func(filter.Filter) filter.Factory
	}{
		{
			name: "Should find filter when exist",
			args: args{
				name: "filter-found",
			},
			want: func(name string) filter.Filter {
				f := &filterMock.Filter{}
				f.On("GetName", mock.Anything).Return(name)
				return f
			},
			getFilterJSONArray: func() []map[string]interface{} {
				return generateFilterAsJSONArray(1)
			},
			getFactory: func(ret filter.Filter) filter.Factory {
				factoryMock := &filterMock.Factory{}
				factoryMock.On("Build", mock.Anything, mock.Anything).Return(ret, nil)
				return factoryMock
			},
		},
		{
			name: "Shoul return error when filters are nil",
			want: func(name string) filter.Filter {
				return nil
			},
			getFilterJSONArray: func() []map[string]interface{} {
				return nil
			},
			getFactory: func(ret filter.Filter) filter.Factory {
				factoryMock := &filterMock.Factory{}
				factoryMock.On("Build", mock.Anything, mock.Anything).Return(nil, errors.New("Error"))
				return factoryMock
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args.name)
			factoryMock := tt.getFactory(want)
			filterArray := tt.getFilterJSONArray()
			d := filter.NewService(factoryMock, filterArray, &kubeMock.Kube{})
			got, err := d.GetFilterByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("dal.GetFilterByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("dal.GetFilterByName() = %v, want %v", got, want)
			}
		})
	}
}
