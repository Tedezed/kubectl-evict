package main

import (
	"fmt"
	"os"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"./lib"
)

const (
	EvictionKind = "Eviction"
)

var (
	namespace = ""
	verbose   = false
)

type (
	ArgOutput = parse_args.ArgOutput
	ArgStruct = parse_args.ArgStruct
)

func main() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	currentNamespace, _, err := kubeConfig.Namespace()
	// if err != nil {
	// 	panic(err.Error())
	// }

	var command ArgOutput
	format_namespace := []string {"--namespace", "-n"}
	fromat_verbose := []string {"--verbose", "-v"}
	var to_find = []ArgStruct {
		ArgStruct {
			Name: "namespace",
			Format: format_namespace,
			GetNext: true,
		},
		ArgStruct {
			Name: "verbose",
			Format: fromat_verbose,
			GetNext: false,
		},
	}

    argsWithoutProg := os.Args[1:]
    command = parse_args.ParseArgs(argsWithoutProg, to_find)

    if (command.ArgMap["namespace"] != "") {
    	namespace = command.ArgMap["namespace"]
    } else {
    	namespace = currentNamespace
    }

    if (command.ArgMap["verbose"] != "") {
    	verbose = true
    }
    //fmt.Println(command.Rest[0], command.ArgMap["namespace"], command.ArgMap["verbose"])



	//flag.StringVar(&namespace, "namespace", currentNamespace, "namespace of the pod")
	//flag.BoolVar(&verbose, "verbose", false, "show more details")
	//flag.Parse()

	//if len(flag.Args()) != 1 {
	if (command.ArgMap["namespace"] == "") {
		fmt.Println("USAGE: kubectl evict POD_NAME [--namespace NAMESPACE]")
		os.Exit(1)
	}
	//podName := flag.Args()[0]
	podName := command.Rest[0]

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		panic(err.Error())
	}

	if verbose {
		fmt.Println("Pod Name:", podName)
		fmt.Println("Namespace:", namespace)
		fmt.Println("Context:", rawConfig.CurrentContext)
	}

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	policyGroupVersion := "v1beta1"
	deleteOptions := &metav1.DeleteOptions{}

	eviction := &policyv1beta1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		DeleteOptions: deleteOptions,
	}
	err = clientset.PolicyV1beta1().Evictions(eviction.Namespace).Evict(eviction)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Done!")
	}

	os.Exit(0)
}
