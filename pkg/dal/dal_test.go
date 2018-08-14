package dal

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetDal(t *testing.T) {
	tests := []struct {
		name string
		want *Dal
	}{
		{
			name: "Should get",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GenerateReasonFilter() *ReasonFilter {
	return &ReasonFilter{
		BaseFilter: BaseFilter{
			Name: "OnlyKillingPods",
			Type: "reason",
		},
		Reason: "Killing",
	}
}

func TestNewDalFromFilePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Dal
	}{
		{
			name: "Should read yaml with filters from testdata/filters.yaml",
			args: args{
				path: "filters.yaml",
			},
			want: &Dal{
				Filters:      []Ifilter{GenerateReasonFilter()},
				Destinations: []IDestination{},
				Integrations: []*Integration{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testdata", tt.args.path)
			filters := NewDalFromFilePath(path).Filters
			wantFilters := tt.want.Filters
			for i := range filters {
				filter := filters[i]
				wantFilter := wantFilters[i]
				name := filter.GetName()
				wantName := wantFilter.GetName()
				if name != wantName {
					t.Errorf("NewDalFromFilePath() = %s, want: %s", name, wantName)
				}
			}
		})
	}
}
