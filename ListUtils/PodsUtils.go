package listutils

import (
	"bytes"
	"context"
	"io"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func GetPodLog(namespace string, podName string, clientset *kubernetes.Clientset) string {
	req := clientset.CoreV1().Pods("monitoring").GetLogs("grafana",&v1.PodLogOptions{
		Previous: true,
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