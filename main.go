package main

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	noProxyForOldMen := os.Getenv("NO_PROXY")
	if noProxyForOldMen != "" {
		ips, err := net.LookupIP(noProxyForOldMen)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		log.Printf("no_proxy=%v ips=%v\n", noProxyForOldMen, ips)
		if len(ips) == 1 {
			os.Setenv("NO_PROXY", ips[0].String())
		}
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	serverVersion, err := clientSet.ServerVersion()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Printf("%#v\n", serverVersion)

	pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("err=`%v`\n", err)
		os.Exit(1)
	}
	log.Printf("pod.count=%d\n", len(pods.Items))

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
