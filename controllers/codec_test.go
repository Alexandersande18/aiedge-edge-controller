package controllers

import (
	// "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	// "k8s.io/apimachinery/pkg/runtime"
	glog "log"
	// "strings"
	// "io/ioutil"
	// "context"
)

func TestDecode(t *testing.T) {
	// c, _ := deserializeAndRenderFromFile("../yamls/test.yaml", "192.168.20.221", "edge1", "registry.aiedge.ndsl-lab.cn")
	// c, err := getAllBaseConfig("../yamls/base/", "192.168.20.221", "registry.aiedge.ndsl-lab.cn")
	// c, err := getAllCloudServiceConfig("../yamls/cloud-app/", "192.168.20.221", "registry.aiedge.ndsl-lab.cn")
	c, err := getAllDeviceModels("../yamls/devicemodel/", "192.168.20.221", "registry.aiedge.ndsl-lab.cn")
	if err != nil {
		glog.Println(err)
	}
	glog.Println(len(c))
	for _, a := range c {
		glog.Println(a.GetNamespace(), a.GetName(), a.GetObjectKind().GroupVersionKind())
		// glog.Println(a)
	}
	// data, _ := ioutil.ReadFile("../yamls/base/jwt-key.yaml")
	// a := strings.Split(string(data), "---")
	// glog.Println(a[1])
}
