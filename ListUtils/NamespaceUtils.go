package listutils

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)



func GetAllNamespacesName(clientset *kubernetes.Clientset) ([]string, error) {

	allNs :=[]string{}

	namespacesListType, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})


	if err != nil {
		log.Fatal(err.Error())
	}

	namespaces := namespacesListType.Items

	for _, namespace := range namespaces{
		allNs = append(allNs, namespace.Name)
	}

	return allNs, err
}

