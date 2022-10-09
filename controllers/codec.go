package controllers

import (
	"strings"

	"k8s.io/client-go/kubernetes/scheme"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "k8s.io/apimachinery/pkg/runtime"
	"io/ioutil"
	glog "log"

	// "context"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func deserialize(data []byte) (client.Object, error) {
	apiextensionsv1.AddToScheme(scheme.Scheme)
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)
	decoder := scheme.Codecs.UniversalDeserializer()
	runtimeObject, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}
	clientObj := runtimeObject.(client.Object)
	return clientObj, nil
}

func deserializeFromFile(fileName string, nodePortIp string, edgeName string, imageRegistry string) ([]client.Object, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
    var ret []client.Object
    for _, doc := range strings.Split(string(data), "---") {
        var afterRender string
        afterRender = strings.ReplaceAll(doc, "{{AIEDGE_SUBNET_NAME}}", edgeName)
        afterRender = strings.ReplaceAll(afterRender, "{{AIEDGE_NODEPORT_IP}}", nodePortIp)
        afterRender = strings.ReplaceAll(afterRender, "{{AIEDGE_IMAGE_REGISTRY}}", imageRegistry)
        glog.Println(afterRender)
        c, err := deserialize([]byte(afterRender))
        if err != nil {
            return ret, err
        }
        ret = append(ret, c)
    }
	return ret, nil
}
