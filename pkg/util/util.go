package util

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func UnmarshalOrDie(in []byte, out interface{}) {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		fmt.Printf("Failed to unmarshal with error: %s", err.Error())
		os.Exit(1)
	}
}

func MapToObjectOrDie(m map[string]interface{}, o interface{}) {
	b, err := yaml.Marshal(m)
	if err != nil {
		fmt.Printf("Failed to marshal with error: %s", err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(b, o)
	if err != nil {
		fmt.Printf("Failed to unmarshal with error: %s", err.Error())
		os.Exit(1)
	}
}

func EchoError(err error) {
	fmt.Printf("Error: %s\n", err.Error())
}
