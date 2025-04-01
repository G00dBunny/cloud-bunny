package listutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)



const (
	CrashLoop string = "CrashLoopBackOff"
	FailedPod string = "Failed"
	ImagePullBackOff string = "ImagePullBackOff"
	Error string = "Error"

)

func GetAllPodsName(namespaceName []string, clientset *kubernetes.Clientset) ([]string) {

	podListName := []string{}

	for _, namespace := range namespaceName{
		podsListType, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Fatal(err.Error())
		}

		pods := podsListType.Items

		// fmt.Printf("\n")

		// fmt.Printf("========================================== \n")

		// fmt.Printf("NAMESPACE : %s \n", namespace)
		
		// fmt.Printf("\n")

		for _, pod := range pods {
			podListName = append(podListName, pod.Name)
			// fmt.Println(podName)
		}

	}

	return podListName
}

// func GetPod(namespace []string, podName string) []string {
	
// 	for _, namespace := range namespaces{
// 		podsListType, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
// 	}

// }

func GetPodLog(namespace string, podName string, clientset *kubernetes.Clientset) string {

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName,&v1.PodLogOptions{
		SinceTime: &metav1.Time{Time: time.Now().Add(time.Duration(-time.Hour))},
		// Previous: true,
	})

	podLogs, err := req.Stream(context.TODO())

	if err != nil {
		log.Fatal(err.Error())
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf,podLogs)

	if err != nil {
		log.Fatal(err.Error())
	}

	return buf.String()
}

func GetPodList(namespaces []string, clientset *kubernetes.Clientset) ([]v1.Pod, error) {
    var allPods []v1.Pod
    
    for _, namespace := range namespaces {
        podsListType, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            return nil, fmt.Errorf("error listing pods in namespace %s: %v", namespace, err)
        }

        allPods = append(allPods, podsListType.Items...)
    }
    
    return allPods, nil
}

func GetBadPod(namespaces []string, clientset *kubernetes.Clientset) []string {
    badPodNames := []string{}
    podList, _ := GetPodList(namespaces, clientset)
    
    for _, pod := range podList {
        if pod.Status.Phase != v1.PodRunning {
            badPodNames = append(badPodNames, pod.Name)
            continue
        }
        
        for _, status := range pod.Status.ContainerStatuses {
            isBad := false
            
            if status.State.Waiting != nil && (
                status.State.Waiting.Reason == CrashLoop ||
                status.State.Waiting.Reason == Error ||
                status.State.Waiting.Reason == ImagePullBackOff) {
                isBad = true
            }
            
            if status.RestartCount > 5 && pod.Status.StartTime != nil {
                isBad = true
            }
            
            if isBad {
                badPodNames = append(badPodNames, pod.Name)
                break 
            }
        }
    }

	// fmt.Println(badPodNames)
    
    return badPodNames
}

func FindPodNamespace(podName string, namespaces []string, clientset *kubernetes.Clientset) string {
	for _, ns := range namespaces {
		pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			continue
		}
		
		for _, pod := range pods.Items {
			if pod.Name == podName {
				return ns
			}
		}
	}
	return ""
}