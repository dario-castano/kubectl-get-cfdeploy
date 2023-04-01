package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	var deployments *appsv1.DeploymentList

	key_flag := flag.String("k", "", "Key to find in the configmaps")
	namespace_flag := flag.String("n", "", "Namespace to be analyzed")

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("c", filepath.Join(home, ".kube", "config"), "(optional) Absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("c", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	configmap_name := flag.Arg(0)
	nargs := flag.NArg()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deployments, err = clientset.AppsV1().Deployments(*namespace_flag).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	if nargs != 0 {
		fmt.Println("---------------------" + configmap_name + "---------------------")

		cfmap, err := clientset.CoreV1().ConfigMaps(*namespace_flag).Get(context.TODO(), configmap_name, metav1.GetOptions{})

		if err != nil {
			panic(err.Error())
		}

		if *key_flag != "" {
			fmt.Println(*key_flag + ": " + cfmap.Data[*key_flag])
		} else {
			for k, v := range cfmap.Data {
				fmt.Println(k + ": " + v)
			}
		}

		fmt.Println("")

	} else {
		for _, v := range deployments.Items {
			replaced := strings.Replace(v.Name, "-deployment", "", 1)
			replaced += "-configmap"

			fmt.Println("---------------------" + replaced + "---------------------")

			cfmap, err := clientset.CoreV1().ConfigMaps(*namespace_flag).Get(context.TODO(), replaced, metav1.GetOptions{})

			if err != nil {
				panic(err.Error())
			}

			if *key_flag != "" {
				fmt.Println(*key_flag + ": " + cfmap.Data[*key_flag])
			} else {
				for k, v := range cfmap.Data {
					fmt.Println(k + ": " + v)
				}
			}
			fmt.Println("")
		}
	}

}
