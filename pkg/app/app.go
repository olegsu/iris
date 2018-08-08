package app

import (
	"github.com/olegsu/iris/pkg/dal"
	"github.com/olegsu/iris/pkg/server"
)

type ApplicationOptions struct {
	IrisPath       string
	KubeconfigPath string
	InCluster      bool
}

func NewApplicationOptions(irisconfig string, kubeconfig string, incluster bool) *ApplicationOptions {
	return &ApplicationOptions{
		IrisPath:       irisconfig,
		KubeconfigPath: kubeconfig,
		InCluster:      incluster,
	}
}

func CreateApp(config *ApplicationOptions) {
	dal.NewDalFromFilePath(config.IrisPath)
	go dal.GetClientset(config.KubeconfigPath, config.InCluster).StartWatching()
	server.StartServer()
}
