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
	c, err := getJWTConfig("../yamls/jwt/")
	if err != nil {
		glog.Println(err)
	}
	glog.Println(c[0].GetName(), c[1].GetName(), c[2].GetName())
	// data, _ := ioutil.ReadFile("../yamls/base/jwt-key.yaml")
	// a := strings.Split(string(data), "---")
	// glog.Println(a[1])
}
