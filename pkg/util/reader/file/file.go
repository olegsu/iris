package file

import "io/ioutil"

// FileReader
type FileReader interface {
	Read(string) ([]byte, error)
}

type reader struct{}

// Read - reads the actual file from path
func (r *reader) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// NewFileReader - creates file reader
func NewFileReader() FileReader {
	return &reader{}
}

// ProcessFile - execute filereader with path as input
func ProcessFile(r FileReader, path string) ([]byte, error) {
	return r.Read(path)
}
