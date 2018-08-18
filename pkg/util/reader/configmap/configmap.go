package configmap

import (
	"github.com/olegsu/iris/pkg/kube"
)

// ConfigmapReader
type ConfigmapReader interface {
	Read(string, string) ([]byte, error)
}

type reader struct {
	kube kube.Kube
}

// Read - read iris data from kubernetes configmap
func (r *reader) Read(name string, namespace string) ([]byte, error) {
	return r.kube.GetIRISConfigmap(namespace, name)
}

// NewConfigmapReader - creates configmap reader
// obj shoud be castable to kubernetes.Clientset
func NewConfigmapReader(k kube.Kube) ConfigmapReader {
	return &reader{
		kube: k,
	}
}

// ProcessConfigmap - execute ConfigmapReader with configmap name and namespace
func ProcessConfigmap(r ConfigmapReader, name string, namespace string) ([]byte, error) {
	return r.Read(name, namespace)
}
