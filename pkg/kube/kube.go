package kube

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/olegsu/iris/pkg/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type WatchFn func(obj interface{})

type Kube interface {
	Watch(WatchFn)
	ResourceByLabelsExist(interface{}, map[string]string) (bool, error)
}

type kube struct {
	Clientset *kubernetes.Clientset
	logger    logger.Logger
}

func (k *kube) Watch(watchFn WatchFn) {
	k.logger.Debug("Starting kuberntes watcher")
	watchlist := cache.NewListWatchFromClient(k.Clientset.Core().RESTClient(), "events", metav1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Event{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				watchFn(obj)
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func (k *kube) ResourceByLabelsExist(obj interface{}, labels map[string]string) (bool, error) {
	selector := ""
	for k, v := range labels {
		if selector == "" {
			selector = fmt.Sprintf("%s=%s", k, v)
		} else {
			selector = fmt.Sprintf("%s,%s=%s", selector, k, v)
		}
	}
	var ev *v1.Event
	bytes, err := json.Marshal(obj)
	json.Unmarshal(bytes, &ev)
	opt := metav1.ListOptions{
		LabelSelector: selector,
	}
	pods, err := k.Clientset.CoreV1().Pods(ev.InvolvedObject.Namespace).List(opt)
	if err != nil {
		return false, err
	}
	return len(pods.Items) > 0, nil
}

func NewKubeManager(kubeconfig string, incluster bool, logger logger.Logger) Kube {
	k := &kube{
		logger: logger,
	}
	var config *rest.Config
	var err error
	if incluster == true {
		logger.Debug("Running from in cluster")
		config, err = rest.InClusterConfig()
	} else {
		logger.Debug("Connecting to cluster from kubeconfig", "path", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	k.Clientset = cs
	return k
}
