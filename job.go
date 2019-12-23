package main

import (
	"fmt"
	"math/rand"
	"time"

	"encoding/json"
	"github.com/nlopes/slack"
	batchv1 "k8s.io/api/batch/v1"
	yaml "k8s.io/apimachinery/pkg/util/yaml"
)

type JobHandler struct {
	client *slack.Client
}

func (j JobHandler) Execute(url string, commands []string, channel string) (job batchv1.Job, err error) {
	rawFile, err := GetFile(url)
	if err != nil {
		return
	}

	jsonFile, err := yaml.ToJSON(rawFile)
	err = json.Unmarshal(jsonFile, &job)
	if err != nil {
		return
	}

	job.Name = job.Name + "-" + RandString(10)
	for i, container := range job.Spec.Template.Spec.Containers {
		job.Spec.Template.Spec.Containers[i].Command = append(container.Command, commands...)
	}

	if err = createJob(&job); err != nil {
		return
	}

	go j.watchAndNotify(job.Name, channel)
	return
}

func (j JobHandler) watchAndNotify(jobName string, channel string) {
	t := time.NewTicker(time.Duration(20) * time.Second)
	for {
		select {
		case <-t.C:
			job, err := getJob(jobName)
			if err != nil {
				fmt.Println("[ERROR] Quit watching job ", job.Name)
				t.Stop()
				return
			}
			fmt.Println("[INFO] Watch job", job.Name)
			if job.Status.Succeeded >= 1 {
				msg := slack.Attachment{Color: "#36a64f", Text: fmt.Sprintf("Succeed %s Job execution", jobName)}
				j.client.PostMessage(channel, slack.MsgOptionAttachments(msg))
				t.Stop()
				return
			}
			if job.Status.Failed >= 1 {
				msg := slack.Attachment{Color: "#e01e5a", Text: fmt.Sprintf("Failed %s execution", jobName)}
				j.client.PostMessage(channel, slack.MsgOptionAttachments(msg))
				t.Stop()
				return
			}
		}
	}
	t.Stop()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
