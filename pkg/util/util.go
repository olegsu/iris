package util

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type UtilService interface {
	ReadFileOrDie(string) []byte
}

type Util struct {
}

func GetUtil() *Util {
	return &Util{}
}

func (util *Util) ReadFileOrDie(path string) []byte {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("Cannot read file %s, exiting\n", path)
		os.Exit(1)
	}
	return file
}

func (util *Util) UnmarshalOrDie(in []byte, out interface{}) {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		fmt.Printf("Failed to unmarshal with error: %s", err.Error())
		os.Exit(1)
	}
}

func (util *Util) MapToObjectOrDie(m map[string]interface{}, o interface{}) {
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
