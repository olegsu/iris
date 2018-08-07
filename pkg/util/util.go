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
