package file

import (
	"io/ioutil"

	"github.com/olegsu/iris/pkg/kube"
)

// FileReader
type FileReader interface {
	Read(string) ([]byte, error)
}

type reader struct {
	kube kube.Kube
}

// Read - reads the actual file from path
func (r *reader) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// NewFileReader - creates file reader
func NewFileReader(k kube.Kube) FileReader {
	return &reader{
		kube: k,
	}
}

// ProcessFile - execute filereader with path as input
func ProcessFile(r FileReader, path string) ([]byte, error) {
	return r.Read(path)
}
