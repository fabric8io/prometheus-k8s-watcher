package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

const (
	nodesTemplateString = `- targets:
{{range .Nodes}}  - {{$address := index .Status.Addresses 0}} {{$address.Address}}:{{$.NodePort}}
{{end}}`
)

var (
	argMaster           = flag.String("master", "https://kubernetes.default.svc.cluster.local", "The URL of the Kubernetes API server")
	argApiVersion       = flag.String("api-version", "v1", "API version to use")
	argInsecure         = flag.Bool("insecure", false, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure.")
	argBearerTokenFile  = flag.String("bearer-token-file", "/var/run/secrets/kubernetes.io/serviceaccount/token", "The file containing the bearer token.")
	argCaCertFile       = flag.String("ca-cert-file", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", "The file containing the CA certificate.")
	argNodeTargetsFile  = flag.String("nodes-file", "/etc/prometheus/config.d/nodes.yml", "The file to write the node targets to.")
	argNodeReadOnlyPort = flag.Int("node-read-only-port", 10255, "The port that metrics can be retrieved from the nodes.")
	nodesTemplate       = template.Must(template.New("nodes").Parse(nodesTemplateString))
)

func main() {
	flag.Parse()

	bearerToken, err := ioutil.ReadFile(*argBearerTokenFile)
	if err != nil {
		log.Fatal(err)
	}

	config := client.Config{
		Host:        *argMaster,
		Insecure:    *argInsecure,
		BearerToken: string(bearerToken),
		Version:     *argApiVersion,
		TLSClientConfig: client.TLSClientConfig{
			CAFile: *argCaCertFile,
		},
	}

	client, err := client.New(&config)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	go watchNodes(client)
	<-done
}

func watchNodes(client *client.Client) {
	nodeList, err := client.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		log.Fatal(err)
	}
	nodes := nodeList.Items
	writeNodeTargetsFile(nodes)
	watcher, err := client.Nodes().Watch(labels.Everything(), fields.Everything(), nodeList.ResourceVersion)
	if err != nil {
		log.Fatal(err)
	}

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added:
			switch obj := event.Object.(type) {
			case *api.Node:
				nodes = append(nodes, *obj)
			}
			writeNodeTargetsFile(nodes)
		case watch.Deleted:
			switch obj := event.Object.(type) {
			case *api.Node:
				index := findNodeIndexInSlice(nodes, obj)
				nodes = append(nodes[:index], nodes[index+1:]...)
			}
			writeNodeTargetsFile(nodes)
		}
	}
}

func findNodeIndexInSlice(slice []api.Node, obj *api.Node) int {
	for i, sliceObj := range slice {
		if sliceObj.ObjectMeta.Namespace == obj.ObjectMeta.Namespace && sliceObj.ObjectMeta.Name == obj.ObjectMeta.Name {
			return i
		}
	}
	return -1
}

func writeNodeTargetsFile(nodes []api.Node) {
	file, err := os.Create(*argNodeTargetsFile)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Nodes    []api.Node
		NodePort int
	}{
		nodes,
		*argNodeReadOnlyPort,
	}
	nodesTemplate.Execute(file, &data)
	file.Close()
}
