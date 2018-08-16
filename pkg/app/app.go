package app

import (
	"github.com/olegsu/iris/pkg/dal"
	"github.com/olegsu/iris/pkg/reader"
	"github.com/olegsu/iris/pkg/server"
)

type ApplicationOptions struct {
	IrisPath               string
	KubeconfigPath         string
	InCluster              bool
	IrisConfigMapName      string
	IrisConfigMapNamespace string
}

func NewApplicationOptions(irisconfig string, kubeconfig string, incluster bool, irisCmName string, irisCmNamespace string) *ApplicationOptions {
	return &ApplicationOptions{
		IrisPath:               irisconfig,
		KubeconfigPath:         kubeconfig,
		InCluster:              incluster,
		IrisConfigMapName:      irisCmName,
		IrisConfigMapNamespace: irisCmNamespace,
	}
}

func CreateApp(config *ApplicationOptions) {
	cs := dal.GetClientset(config.KubeconfigPath, config.InCluster)
	var r reader.IRISProcessor
	if config.IrisConfigMapName != "" {
		r, _ = reader.NewProcessor([]string{
			config.IrisConfigMapName,
			config.IrisConfigMapNamespace,
		}, nil)
	} else {
		r, _ = reader.NewProcessor([]string{
			config.IrisPath,
		}, cs.Clientset.CoreV1())
	}
	bytes, _ := reader.Process(r)
	dal.CreateDalFromBytes(bytes)
	go cs.StartWatching()
	server.StartServer()
}
