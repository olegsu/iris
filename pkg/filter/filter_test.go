package filter_test

import (
	"testing"

	"github.com/olegsu/iris/pkg/filter"
	"github.com/olegsu/iris/pkg/filter/mocks"
	"github.com/stretchr/testify/mock"
)

func TestApplyFilter(t *testing.T) {
	type args struct {
		f   func() filter.Filter
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should call apply on filter",
			args: args{
				f: func() filter.Filter {
					m := &mocks.Filter{}
					m.On("Apply", mock.Anything).Return(true, nil)
					return m
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filter.ApplyFilter(tt.args.f(), tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ApplyFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GenerateJSON() map[string]interface{} {
	return map[string]interface{}{
		"name": "name",
		"type": "jsonpath",
	}
}

// func Test_baseFilter(t *testing.T) {

// 	tests := []struct {
// 		name      string
// 		getFilter func() filter.Filter
// 		wantName  string
// 		wantType  string
// 	}{
// 		{
// 			name: "Get filter name",
// 			getFilter: func() filter.Filter {
// 				factory := filter.NewFactory()
// 				f, _ := factory.Build(GenerateJSON(), &kube.Kube{})
// 				return f
// 			},
// 			wantName: "name",
// 			wantType: "jsonpath",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := tt.getFilter()
// 			if got := f.GetName(); got != tt.wantName {
// 				t.Errorf("GetName() = %v, want %v", got, tt.wantName)
// 			}

// 			if got := f.GetType(); got != tt.wantType {
// 				t.Errorf("GetName() = %v, want %v", got, tt.wantType)
// 			}
// 		})
// 	}
// }

func Test_baseFilter_GetType(t *testing.T) {

	tests := []struct {
		name      string
		getFilter func() filter.Filter
		want      string
	}{
		{
			name: "Get filter name",
			getFilter: func() filter.Filter {
				f := &mocks.Filter{}
				f.On("GetType", mock.Anything).Return("filter-type")
				return f
			},
			want: "filter-type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.getFilter()
			if got := f.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
