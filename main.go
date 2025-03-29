/*
*   IDEAS
*			- use lightweight kubeconfig context (long-lived token if in aws)
*			- Skip pods in Running state with no recent restarts
*			- Only fetch logs -> error, CrashLoopBackOff
*			- if in aws : run in small EC2 instance inside the same VPC as eks cluster?
*			- For logs : use TailLines + SinceSeconds to get last numbers of lines or last seconds
*			- track restart counts of pods -> detect pods with high restarts
*			- mark pods as stale
*			- integrate  pod status and health
*			- memory estimate
*			-
*	TODO :
*			- In gargamel : add flush to autoremove expired entries
*			- background go routine that purges old pods
*			-
*			-
*			-
*
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/G00dBunny/Gargamel/gargamel"
	"github.com/G00dBunny/cloud-bunny/listutils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
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
	clientset.CoreV1()

	if err != nil {
		log.Fatal(err.Error())
	}

	exp := gargamel.Expiration(gargamel.NoExpiration)


	cache := gargamel.New(&exp)

	ns := "monitoring"

	namespace := gargamel.Namespace{
		Name: ns,
	}
	logs := listutils.GetPodLog("monitoring", "real-memory-leak", clientset)

	podList := []*gargamel.Pod{
		{Name: logs},
	}

	cache.Set(&namespace, podList)
	

	fmt.Println(cache.String())
	

	namespaceList,_ := listutils.GetAllNamespacesName(clientset)

	listutils.GetBadPod(namespaceList,clientset)
	

}
