/*
  ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
  │                                                                                                                 │
  │ *  DONE                                                                                                         │
  │ *      - Only fetch logs -> error, CrashLoopBackOff                                                             │
  │ *      - For logs : use TailLines + SinceSeconds to get last numbers of lines or last seconds                   │
  │ *      - Skip pods in Running state with no recent restarts                                                     │
  │ *      - integrate  pod status and health                                                                       │
  │ *                                                                                                               │
  │ *   NOTE                                                                                                       │
  │ *      - use lightweight kubeconfig context (long-lived token if in aws)                                        │
  │ *      - if in aws : run in small EC2 instance inside the same VPC as eks cluster?                              │
  │ *      - track restart counts of pods -> detect pods with high restarts                                         │
  │ *      - mark pods as stale                                                                                     │
  │ *      - memory estimate                                                                                        │
  │ *      -                                                                                                        │
  │ *      - Automated jira ticket creation after vulnerability detected -> have the use decide if it wants to      │
  │ create it or                                                                                                    │
  │ not                                                                                                             │
  │ *      -                                                                                                        │
  │ *      -                                                                                                        │
  │ *      -                                                                                                        │
  │ *  TODO :                                                                                                       │
  │ *      - In gargamel : add flush to autoremove expired entries                                                  │
  │ *      - background go routine that purges old pods                                                             │
  │ *      -                                                                                                        │
  │ *      -                                                                                                        │
  │ *      -                                                                                                        │
  │ *                                                                                                               │
  │                                                                                                                 │
  └─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
*/

/*
┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ // ANCHOR - Used to indicate a section in your file                                                     │
│ // TODO - An item that is awaiting completion                                                           │
│ // FIXME - An item that requires a bugfix                                                               │
│ // STUB - Used for generated default snippets                                                           │
│ // NOTE - An important note for a specific code section                                                 │
│ // REVIEW - An item that requires additional review                                                     │
│ // SECTION - Used to define a region (See 'Hierarchical anchors')                                       │
│ // DONE - Used for Done issues 															              │
│                                                                                                         │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/G00dBunny/cloud-bunny/config"
	"github.com/G00dBunny/cloud-bunny/jiraBed"
	"github.com/G00dBunny/cloud-bunny/listutils"
	"github.com/G00dBunny/cloud-bunny/llmBed"
	"github.com/G00dBunny/cloud-bunny/prompt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	config.LoadEnv()
	
	apiKey, model, maxTokens, timeoutSec := config.GetOpenAIConfig()
	if apiKey == "" {
		log.Println("OPENAI_API_KEY not set in environment")
	}


	
	jiraConfig, err := config.GetJiraConfig()


	jiraConfigured := jiraConfig.Username != "" && jiraConfig.Token != "" && jiraConfig.URL != "" && jiraConfig.Project != ""

	// createTicket := jiraConfig.IsUsed

	if err != nil {
		log.Println("Jira configuration incomplete : ")
	}

	
	
	/*
	* NOTE : https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
	*/
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}


	clientset, err := kubernetes.NewForConfig(config)
	clientset.CoreV1()

	if err != nil {
		log.Fatal(err.Error())
	}

	namespaceList,_ := listutils.GetAllNamespacesName(clientset)

	badPods := listutils.GetBadPod(namespaceList, clientset)
	if len(badPods) == 0 {
		fmt.Println("No bad pods detected.")
		return
	}
	
	fmt.Printf("Found %d bad pods: %v\n", len(badPods), badPods)

	var analysis []llmBed.AnalysisResult


	if apiKey != "" {
		fmt.Println("\nStarting analysis of problematic pods...")
		
		gptConfig := llmBed.NewConfig(apiKey, model, maxTokens, timeoutSec)
		
		analysis = llmBed.AnalyzeBadPods(namespaceList, clientset, gptConfig)
		
		for _, result := range analysis {
			fmt.Printf("\n=============================================\n")
			fmt.Printf("POD: %s (Namespace: %s)\n", result.PodName, result.Namespace)
			fmt.Printf("=============================================\n")
			
			if result.Error != nil {
				fmt.Printf("Error: %v\n", result.Error)
				continue
			}
			
			fmt.Printf("ANALYSIS:\n%s\n", result.Analysis)
			
			if jiraConfigured {
				fmt.Printf("\nDo you want to create a Jira ticket for this issue? (y/n): ")
				var response string
				fmt.Scanln(&response)
				
				if response == "y" || response == "Y" {
					ticketDetail := jiraBed.GenerateTicketFromAnalysis(result)
					
					
					if prompt.PromptUserForTicketConfirmation(&ticketDetail) {
					
						jiraClient, err := jiraBed.JiraClient(jiraConfig)
						if err != nil {
							log.Printf("Jira Client failed: %v", err)
							continue
						}
						
						issue, err := jiraBed.CreateTicket(jiraClient, jiraConfig.Project, ticketDetail)
						if err != nil {
							log.Printf("Failed to create ticket for pod %s in namespace %s: %v", 
								result.PodName, result.Namespace, err)
							continue
						}
						
						fmt.Printf("Created Jira ticket %s for pod %s in namespace %s\n", 
							issue.Key, result.PodName, result.Namespace)
						fmt.Printf("URL: %s/browse/%s\n", jiraConfig.URL, issue.Key)
					} else {
						fmt.Println("Ticket creation canceled by user.")
					}
				}
			}
		}
	} else {
		fmt.Println("Skipping analysis due to missing OpenAI API key.")
	}


	// /*
	// *	FIXME : REMOVE THIS PART WHEN DONE WITH TICKET AND GPT IMPLEMENTATION 
	// */
	// if createTicket && jiraConfigured && len(analysis) > 0 {
	// 	fmt.Println("\nCreating Jira tickets...")
		
	// 	config := jiraConfig
		
	// 	jiraClient, err := jiraBed.JiraClient(config)
	// 	if err != nil {
	// 		log.Printf("Jira Client failed: %v", err)
	// 		return
	// 	}
		
	// 	for _, result := range analysis {
	// 		if result.Error != nil {
	// 			log.Printf("Skipping ticket creation for pod %s in namespace %s due to analysis error", 
	// 				result.PodName, result.Namespace)
	// 			continue
	// 		}
			
	// 		ticketDetail := jiraBed.GenerateTicketFromAnalysis(result)
			
	// 		fmt.Printf("Creating ticket for pod: %s in namespace: %s\n", result.PodName, result.Namespace)
	// 		issue, err := jiraBed.CreateTicket(jiraClient, jiraConfig.Project, ticketDetail)
	// 		if err != nil {
	// 			log.Printf("Failed to create ticket for pod %s in namespace %s: %v", 
	// 				result.PodName, result.Namespace, err)
	// 			continue
	// 		}
			
	// 		fmt.Printf("Created Jira ticket %s for pod %s in namespace %s\n", 
	// 			issue.Key, result.PodName, result.Namespace)
	// 	}
	// }

}