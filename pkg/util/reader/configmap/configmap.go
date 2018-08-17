package configmap

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ConfigmapReader
type ConfigmapReader interface {
	Read(string, string) ([]byte, error)
}

type reader struct {
	kube kubernetes.Clientset
}

// Read - read iris data from kubernetes configmap
func (r *reader) Read(name string, namespace string) ([]byte, error) {
	cm, err := r.kube.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return []byte(cm.Data["iris"]), nil
}

// NewConfigmapReader - creates configmap reader
// obj shoud be castable to kubernetes.Clientset
func NewConfigmapReader(obj interface{}) ConfigmapReader {
	kube, ok := obj.(kubernetes.Clientset)
	if ok == false {
		return &reader{}
	}
	return &reader{
		kube: kube,
	}
}

// ProcessConfigmap - execute ConfigmapReader with configmap name and namespace
func ProcessConfigmap(r ConfigmapReader, name string, namespace string) ([]byte, error) {
	return r.Read(name, namespace)
}
