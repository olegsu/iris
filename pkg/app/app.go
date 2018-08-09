package app

import (
	"github.com/olegsu/iris/pkg/dal"
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
	if config.IrisConfigMapName != "" {
		dal.NewDalFromConfigMap(config.IrisConfigMapName, config.IrisConfigMapNamespace)
	} else {
		dal.NewDalFromFilePath(config.IrisPath)
	}
	go cs.StartWatching()
	server.StartServer()
}
