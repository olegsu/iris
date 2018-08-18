package reader

import (
	"reflect"
	"testing"

	"github.com/olegsu/iris/pkg/kube"

	"github.com/olegsu/iris/pkg/util/reader/configmap"
	"github.com/olegsu/iris/pkg/util/reader/file"
)

type mockFileReader struct{}

func (m *mockFileReader) Read(path string) ([]byte, error) {
	return []byte{}, nil
}

type mockConfigmapReader struct{}

func (m *mockConfigmapReader) Read(name string, namespace string) ([]byte, error) {
	return []byte{}, nil
}

type mockKube struct{}

func (m *mockKube) Watch(fn kube.WatchFn) {

}
func (m *mockKube) GetIRISConfigmap(name string, ns string) ([]byte, error) {
	return []byte{}, nil
}
func (m *mockKube) FindResourceByLabels(obj interface{}, labels map[string]string) (bool, error) {
	return true, nil
}

func Test_processor_Process(t *testing.T) {

	type fields struct {
		fileReader      file.FileReader
		configmapReader configmap.ConfigmapReader
		args            []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Should process from file",
			fields: fields{
				fileReader: &mockFileReader{},
				args: []string{
					"path/to/file",
				},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Should process from configmap",
			fields: fields{
				args: []string{
					"name",
					"namespace",
				},
				configmapReader: &mockConfigmapReader{},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Should return error when no readers found ",
			fields: fields{
				args: []string{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &processor{
				fileReader:      tt.fields.fileReader,
				configmapReader: tt.fields.configmapReader,
				args:            tt.fields.args,
			}
			got, err := i.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("processor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processor.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProcessor(t *testing.T) {
	type args struct {
		args []string
		obj  kube.Kube
	}
	tests := []struct {
		name    string
		args    args
		want    IRISProcessor
		wantErr bool
	}{
		{
			name: "Should create processor with file reader when given input len=1",
			args: args{
				args: []string{"path/to/file"},
				obj:  &mockKube{},
			},
			want: &processor{
				args:       []string{"path/to/file"},
				fileReader: file.NewFileReader(&mockKube{}),
			},
			wantErr: false,
		},
		{
			name: "Should create processor with configmap reader when given input len=2",
			args: args{
				args: []string{
					"name",
					"namespace",
				},
				obj: &mockKube{},
			},
			want: &processor{
				args: []string{
					"name",
					"namespace",
				},
				configmapReader: configmap.NewConfigmapReader(&mockKube{}),
			},
			wantErr: false,
		},
		{
			name: "Should return an error when args len is not 1 or 2",
			args: args{
				args: []string{},
				obj:  &mockKube{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProcessor(tt.args.args, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockIRISProcessor struct{}

func (m *mockIRISProcessor) Process() ([]byte, error) {
	return []byte{}, nil
}

func TestProcess(t *testing.T) {
	type args struct {
		processor IRISProcessor
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Should call process on given processor",
			args: args{
				processor: &mockIRISProcessor{},
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Process(tt.args.processor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
