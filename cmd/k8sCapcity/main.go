package main

import (
	"flag"
	"os"
	"path/filepath"

	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
	// Support gcp and other authentication schemes
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	nodeLabel := flag.String("nodelabel", "", "Label to match for nodes, if blank grab all nodes")
	nameSpace := flag.String("namespace", "", "Namespace to grab capacity usage from")
	daemonMode := flag.Bool("daemon", false, "Run in daemon mode")
	jsonMode := flag.Bool("json", false, "Output information in json format")
	checkMode := flag.Bool("check", false, "Check kubernetes connection")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// no config, maybe we are inside a kubernetes cluster.
		config, err = rest.InClusterConfig()
		check(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	check(err)

	if *checkMode {
		_, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		check(err)
		fmt.Println("ok")
		return
	}

	// BreakOut to namespace if asked
	if *nameSpace != "" {
		nsInfo := gatherNamespaceInfo(clientset, nameSpace)
		if *jsonMode {
			result, err := json.Marshal(nsInfo)
			check(err)
			fmt.Println(string(result))
			return
		}
		output := namespaceHumanMode(nsInfo)
		for _, line := range output {
			fmt.Println(line)
		}
		return
	}

	// Gather info
	if *daemonMode {
		for {
			clusterInfo := gatherInfo(clientset, nodeLabel)
			getCapcity(clusterInfo)
			time.Sleep(300 * time.Second)
		}
	} else if *jsonMode {
		clusterInfo := gatherInfo(clientset, nodeLabel)
		getCapcity(clusterInfo)

	} else {
		clusterInfo := gatherInfo(clientset, nodeLabel)
		humanMode(clusterInfo)
	}
}
