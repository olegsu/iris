package reader

import (
	"fmt"

	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/util/reader/file"
)

type (
	// IRISProcessor knows how to process iris configs
	IRISProcessor interface {
		Process() ([]byte, error)
	}

	processor struct {
		fileReader file.FileReader
		args       []string
	}
)

// Process - start processing the iris confi
func (i *processor) Process() ([]byte, error) {
	if i.fileReader != nil {
		return i.fileReader.Read(i.args[0])
	} else {
		return nil, fmt.Errorf("No reader found")
	}
}

// NewProcessor - crete new processor based on the len of the args
func NewProcessor(args []string, k kube.Kube) (IRISProcessor, error) {
	if len(args) == 1 {
		return &processor{
			fileReader: file.NewFileReader(k),
			args:       args,
		}, nil
	} else {
		return nil, fmt.Errorf("Could not create iris processor, arguments not match")
	}
}

// Process - execte processor
func Process(processor IRISProcessor) ([]byte, error) {
	return processor.Process()
}
