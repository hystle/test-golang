package testKubePkg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube struct {
	connected bool
	client    *kubernetes.Clientset
}

type KubeWatchCtx struct {
	ns       string
	resName  string
	resType  string
	watcher  watch.Interface
	stopChan chan bool
}

func StartKubeResWatch(kube *Kube, ns string, resName string, resType string) *KubeWatchCtx {
	if !kube.connected {
		fmt.Println("K8s client is not connected")
		return nil
	}

	watchCtx := new(KubeWatchCtx)
	watchCtx.ns = ns
	watchCtx.resName = resName
	watchCtx.resType = resType
	watchCtx.stopChan = make(chan bool)

	var selector string
	if resName != "" {
		selector = fields.OneTermEqualSelector("metadata.name", resName).String()
	}

	timeout := int64(60 * 60 * 24) // 24 hours
	var watch watch.Interface
	var err error
	if resType == "configmap" {
		listOptions := metav1.ListOptions{FieldSelector: selector, Watch: true, TimeoutSeconds: &timeout}
		watch, err = kube.client.CoreV1().ConfigMaps(ns).Watch(context.TODO(), listOptions)
	} else {
		fmt.Printf("Resource type %s not supported yet", resType)
		return nil
	}
	if err != nil {
		fmt.Printf("Cannot create watch for the %s %s", resType, resName)
		return nil
	}
	watchCtx.watcher = watch

	go kube.watchResRoutine(watchCtx)

	return watchCtx
}

func (kube *Kube) watchResRoutine(watchCtx *KubeWatchCtx) {
ExitRoutine:
	for {
		fmt.Println("-- Watching...")
		// blocks until one of its cases can run
		select {
		case event, ok := <-watchCtx.watcher.ResultChan():
			fmt.Printf("Received event type: %s\n", event.Type)
			if !ok {
				fmt.Printf("Watch Chan for %s closed", watchCtx.resName)
				break ExitRoutine
			}

			// type assertion test
			obj, ok := event.Object.(*corev1.ConfigMap)
			if !ok {
				fmt.Println("Bad type - not a core v1 ConfigMap")
			}
			fmt.Printf("ConfigMap: %s - %s\n", obj.Name, obj.Namespace)
			for k, v := range obj.Data {
				fmt.Printf("Key: %s => Value: %s\n", k, v)
			}
		case <-watchCtx.stopChan:
			fmt.Println("Got terminate signal")
			break ExitRoutine
		}
	}
}

func GetPods(kube *Kube) {
	for {
		// query all pods in ns
		ns := "kube-system"
		pods, err := kube.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})

		// use helper function for error handling
		if errors.IsNotFound(err) {
			fmt.Printf("No Pods in %s namespace\n", ns)
			// Cast to StatusError and use its properties. ex. ErrStatus.Message
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting Pods %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			for _, pod := range pods.Items {
				fmt.Printf("Pods in the cluster: %s, %s\n", pod.Name, pod.Spec.NodeName)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (k *Kube) Connect() bool {
	if k.connected {
		return true
	}

	// 1. create config for client
	var config *rest.Config
	var err error

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	_, err = os.Stat(kubeconfig)
	if err == nil {
		fmt.Println("Use out-of-cluster config")
		// use kubeconfig file if it exists
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		fmt.Println("Use in-cluster config")
		// use Service Account token mounted inside the Pod at the /var/run/secrets/kubernetes.io/serviceaccount
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Printf("Error creating config: %v", err.Error())
		return false
	}

	// 2. create a clientset for the given config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating kube client: %v", err.Error())
		return false
	}

	k.client = clientset
	k.connected = true
	return true
}
