package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/G00dBunny/cloud-bunny/ListUtils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	ListUtils.ListCluster()
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatal(err.Error())
	}


	namespacesListType, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})


	if err != nil {
		log.Fatal(err.Error())
	}

	namespaces := namespacesListType.Items


	for _, namespace := range namespaces{
		nameNS := namespace.Name
		podsListType, err := clientset.CoreV1().Pods(nameNS).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Fatal(err.Error())
		}

		pods := podsListType.Items

		fmt.Printf("\n")

		fmt.Printf("========================================== \n")

		fmt.Printf("NAMESPACE : %s \n", nameNS)
		
		fmt.Printf("\n")

		for _, pod := range pods {
			podName := pod.Name
			fmt.Println(podName)
		}

	}

	
	

}
