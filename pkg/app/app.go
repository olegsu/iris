package app

import (
	"github.com/olegsu/iris/pkg/dal"
	"github.com/olegsu/iris/pkg/kube"
	"github.com/olegsu/iris/pkg/reader"
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
	k := kube.NewKubeManager(config.KubeconfigPath, config.InCluster)
	var r reader.IRISProcessor

	r, _ = reader.NewProcessor([]string{
		config.IrisPath,
	}, k)
	bytes, _ := reader.Process(r)
	d := dal.CreateDalFromBytes(bytes, k)
	fn := func(obj interface{}) {
		onAdd(d.Integrations, obj)
	}
	go k.Watch(fn)
	server.StartServer()
}

func onAdd(integrations []*dal.Integration, obj interface{}) {
	for _, integration := range integrations {
		go integration.Exec(obj)
	}
}
