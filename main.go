package main

import (
	"bytes"
	"regexp"
	"strconv"

	"log"

	"github.com/Albertwzp/cli-go/config"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func main() {
	k, _ := config.GetK8sConfig()
	clientSet := kubernetes.NewForConfigOrDie(k)
	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, 0)
	cmInformer := sharedInformerFactory.Core().V1().ConfigMaps()
	informer := cmInformer.Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	sharedInformerFactory.Start(stopCh)
	sharedInformerFactory.WaitForCacheSync(stopCh)
	/*indexer := cmInformer.Lister()
	cms, err := indexer.List(labels.Everything())
	if err != nil {
		log.Panic(err)
	}
	for _, cm := range cms {
		checkRS(cm)
	}*/
	<-stopCh
}

func onAdd(obj interface{}) {
	cm := obj.(*corev1.ConfigMap)
	checkRS(cm)
}
func onUpdate(oldObj interface{}, newObj interface{}) {
	cm := newObj.(*corev1.ConfigMap)
	checkRS(cm)
}

func onDelete(obj interface{}) {
	cm := obj.(*corev1.ConfigMap)
	if _, ok := cm.Data["application.yaml"]; ok {
		log.Println(cm.Name, ":Delete")
	}
}

func checkRS(cm *corev1.ConfigMap) {
	re := regexp.MustCompile("[0-9]+")
	if app, ok := cm.Data["application.yaml"]; ok {
		viper.SetConfigType("YAML")
		if err := viper.ReadConfig(bytes.NewBuffer([]byte(app))); err == nil {
			log.Println(cm.Name, ":Success")
		} else {
			num := re.FindString(err.Error())
			n, _ := strconv.Atoi(num)
			log.Println(cm.Name, ":Error", n+7)
		}
	}
}
