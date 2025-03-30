/*
*	DONE
*			- Only fetch logs -> error, CrashLoopBackOff
*			- For logs : use TailLines + SinceSeconds to get last numbers of lines or last seconds
*			- Skip pods in Running state with no recent restarts
*			- integrate  pod status and health
*
*   IDEAS
*			- use lightweight kubeconfig context (long-lived token if in aws)
*			- if in aws : run in small EC2 instance inside the same VPC as eks cluster?
*			- track restart counts of pods -> detect pods with high restarts
*			- mark pods as stale
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

	"github.com/G00dBunny/cloud-bunny/config"
	"github.com/G00dBunny/cloud-bunny/listutils"
	"github.com/G00dBunny/cloud-bunny/llmBed"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	config.LoadEnv()
	
	apiKey, model, maxTokens, timeoutSec := config.GetOpenAIConfig()
	if apiKey == "" {
		log.Println("OPENAI_API_KEY not set in environment or .env file; will skip GPT analysis")
	}

	
	/*
	* 	From client-go k8 doc
	*/
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

	// exp := gargamel.Expiration(gargamel.NoExpiration)


	// cache := gargamel.New(&exp)

	// ns := "monitoring"

	// namespace := gargamel.Namespace{
	// 	Name: ns,
	// }
	// logs := listutils.GetPodLog("monitoring", "real-memory-leak", clientset)

	// podList := []*gargamel.Pod{
	// 	{Name: logs},
	// }

	// cache.Set(&namespace, podList)
	

	// fmt.Println(cache.String())
	

	namespaceList,_ := listutils.GetAllNamespacesName(clientset)

	badPods := listutils.GetBadPod(namespaceList, clientset)
	if len(badPods) == 0 {
		fmt.Println("No bad pods detected.")
		return
	}
	
	fmt.Printf("Found %d bad pods: %v\n", len(badPods), badPods)

	if apiKey != "" {
		fmt.Println("\nStarting analysis of problematic pods...")
		
		gptConfig := llmBed.NewConfig(apiKey, model, maxTokens, timeoutSec)
		
		results := llmBed.AnalyzeBadPods(namespaceList, clientset, gptConfig)
		
		for _, result := range results {
			fmt.Printf("\n=============================================\n")
			fmt.Printf("POD: %s (Namespace: %s)\n", result.PodName, result.Namespace)
			fmt.Printf("=============================================\n")
			
			if result.Error != nil {
				fmt.Printf("Error: %v\n", result.Error)
				continue
			}
			
			fmt.Printf("ANALYSIS:\n%s\n", result.Analysis)
		}
	} 

}
