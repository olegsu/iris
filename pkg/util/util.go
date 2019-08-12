package util

import (
	"fmt"
	"os"

	"github.com/olegsu/iris/pkg/logger"

	yaml "gopkg.in/yaml.v2"
)

var (
	BuildVersion, BuildDate, GitTag string
)

func UnmarshalOrDie(in []byte, out interface{}, logger logger.Logger) {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		logger.Error("Failed to unmarshal", "err", err.Error())
		os.Exit(1)
	}
}

func MapToObjectOrDie(m map[string]interface{}, o interface{}, logger logger.Logger) {
	b, err := yaml.Marshal(m)
	if err != nil {
		logger.Error("Failed to marshal", "err", err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(b, o)
	if err != nil {
		logger.Error("Failed to unmarshal", "err", err.Error())
		os.Exit(1)
	}
}

func EchoError(err error) {
	fmt.Printf("Error: %s\n", err.Error())
}
