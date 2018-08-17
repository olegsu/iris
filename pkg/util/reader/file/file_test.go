package file

import (
	"fmt"
	"reflect"
	"testing"
)

type mockFileReaderSuccess struct {
	FileReader
}

func (m *mockFileReaderSuccess) Read(string) ([]byte, error) {
	return []byte{}, nil
}

type mockFileReaderNotFound struct {
	FileReader
}

func (m *mockFileReaderNotFound) Read(path string) ([]byte, error) {
	return nil, fmt.Errorf("%s not found", path)
}

func TestProcessFile(t *testing.T) {
	type args struct {
		r    FileReader
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Should return file content",
			args: args{
				r:    &mockFileReaderSuccess{},
				path: "fath/to/file",
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Should return error if file not found",
			args: args{
				r:    &mockFileReaderNotFound{},
				path: "file/not/found",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessFile(tt.args.r, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFileReader(t *testing.T) {
	tests := []struct {
		name string
		want FileReader
	}{
		{
			name: "Should return Filereader interface",
			want: &reader{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileReader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
