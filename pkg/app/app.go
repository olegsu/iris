package app

import (
	"github.com/olegsu/iris/pkg/dal"
	"github.com/olegsu/iris/pkg/server"
)

func CreateApp(irisconfig string, kubeconfig string) {
	dal.NewDalFromFilePath(irisconfig)
	go dal.GetClientset(kubeconfig).StartWatching()
	server.StartServer()
}
