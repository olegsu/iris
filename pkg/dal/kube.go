package dal

import (
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var kube *KubeManager

type KubeManager struct {
	Clientset *kubernetes.Clientset
}

func GetClientset(kubeconfig string, incluster bool) *KubeManager {
	if kube != nil {
		return kube
	}
	kube = &KubeManager{}
	var config *rest.Config
	var err error
	if incluster == true {
		fmt.Printf("Running from in cluster\n")
		config, err = rest.InClusterConfig()
	} else {
		fmt.Printf("Connecting to cluster from kubeconfig %s\n", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	kube.Clientset = cs
	return kube
}

func onAdd(obj interface{}, cnf *Dal) {
	for index := 0; index < len(cnf.Integrations); index++ {
		go (cnf.Integrations[index]).Exec(obj)
	}
}

func (k *KubeManager) StartWatching() {
	fmt.Printf("Watching\n")
	cnf := GetDal()
	watchlist := cache.NewListWatchFromClient(k.Clientset.Core().RESTClient(), "events", metav1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Event{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAdd(obj, cnf)
			},
			DeleteFunc: func(obj interface{}) {
				onAdd(obj, cnf)
			},
			UpdateFunc: func(old interface{}, new interface{}) {
				onAdd(new, cnf)
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func GetConfigmapData(configmapName string, configmapNamespace string) (string, error) {
	cm, err := kube.Clientset.CoreV1().ConfigMaps(configmapNamespace).Get(configmapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	stringAsData := cm.Data["iris"]
	return stringAsData, nil
}

func (k *KubeManager) FindResourceByLabels(obj interface{}, labels map[string]string) (bool, error) {
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
