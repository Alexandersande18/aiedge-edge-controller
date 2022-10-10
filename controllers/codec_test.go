package controllers

import (
	// "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	// "k8s.io/apimachinery/pkg/runtime"
	glog "log"
	// "context"
)

func TestDecode(t *testing.T) {
	c, _ := deserializeAndRenderFromFile("../yamls/test.yaml", "192.168.20.221", "edge1", "registry.aiedge.ndsl-lab.cn")
	glog.Println(c)
}
