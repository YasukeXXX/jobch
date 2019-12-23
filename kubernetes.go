package main

import (
	"log"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createJob(job *batchv1.Job) (err error) {
	client, err := newKubernetesClient()
	if err != nil {
		log.Print(err)
		return
	}

	job, err = client.BatchV1().Jobs("default").Create(job)
	if err != nil {
		log.Print(err)
		return
	}

	return
}

func getJob(jobName string) (job *batchv1.Job, err error) {
	client, err := newKubernetesClient()
	if err != nil {
		log.Print(err)
		return
	}

	job, err = client.BatchV1().Jobs("default").Get(jobName, metav1.GetOptions{})
	if err != nil {
		log.Print(err)
		return
	}

	return
}

func newKubernetesClient() (kubernetes.Interface, error) {
	if os.Getenv("LOCAL") != "" {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
		return kubernetes.NewForConfig(config)
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
