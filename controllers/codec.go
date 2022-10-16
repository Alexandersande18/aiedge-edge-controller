package controllers

import (
	"strings"

	"k8s.io/client-go/kubernetes/scheme"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "k8s.io/apimachinery/pkg/runtime"
	"io/ioutil"
	// glog "log"
	"path/filepath"
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

func deserializeAndRenderFromFile(fileName string, nodePortIp string, edgeName string, imageRegistry string) ([]client.Object, error) {
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
		// glog.Println(afterRender)
		c, err := deserialize([]byte(afterRender))
		if err != nil {
			return ret, err
		}
		ret = append(ret, c)
	}
	return ret, nil
}

func getSecrets(yamlDir string) ([]client.Object, error) {
	yamlList := []string{"jwt-private-key", "jwt-public-key", "registry-pull-secret"}
	nsList := []string{"aiedge", "aiedge-public-device"}
	var ret []client.Object
	for _, ym := range yamlList {
		if ym == "jwt-private-key" {
			data, err := ioutil.ReadFile(filepath.Join(yamlDir, ym+".yaml"))
			c, err := deserialize(data)
			if err != nil {
				return ret, err
			}
			ret = append(ret, c)
			// glog.Println(ret)
		} else {
			for _, ns := range nsList {
				data, err := ioutil.ReadFile(filepath.Join(yamlDir, ym+".yaml"))
				afterRender := strings.ReplaceAll(string(data), "{{AIEDGE_NAMESPACE}}", ns)
				c, err := deserialize([]byte(afterRender))
				if err != nil {
					return ret, err
				}
				ret = append(ret, c)
			}
		}

	}
	return ret, nil
}

func getAllPerEdgeServiceConfig(yamlDir string, nodePortIp string, edgeName string, imageRegistry string, baseArch string) ([]client.Object, error) {
	yamlList := []string{"cece", "edge-stream-relay", "stream-transfer"}
	arch := ""
	if baseArch == "arm64" {
		arch = "-arm64"
	}
	var ret []client.Object
	for _, ym := range yamlList {
		ls, err := deserializeAndRenderFromFile(filepath.Join(yamlDir, ym+arch+".yaml"), nodePortIp, edgeName, imageRegistry)
		if err != nil {
			return ret, err
		}
		for _, obj := range ls {
			ret = append(ret, obj)
		}
	}
	return ret, nil
}

func getAllBaseConfig(yamlDir string, brokerClusterIp string, imageRegistry string) ([]client.Object, error) {
	yamlList := []string{"coredns", "mysql-pv", "rocketmq-acl-conf", "mysql-deploy", "rocketmq-cluster-x", "stream-srs"}
	var ret []client.Object
	for _, ym := range yamlList {
		ls, err := deserializeAndRenderFromFile(filepath.Join(yamlDir, ym+".yaml"), "", "", imageRegistry)
		if err != nil {
			return ret, err
		}
		for _, obj := range ls {
			ret = append(ret, obj)
		}
	}
	return ret, nil
}

func getAllCloudServiceConfig(yamlDir string, schedulerNodePortIp string, imageRegistry string) ([]client.Object, error) {
	yamlList := []string{"auth-deploy", "acc", "acc-front", "cec-cluster", "stream-scheduler"}
	var ret []client.Object
	for _, ym := range yamlList {
		ls, err := deserializeAndRenderFromFile(filepath.Join(yamlDir, ym+".yaml"), schedulerNodePortIp, "", imageRegistry)
		if err != nil {
			return ret, err
		}
		for _, obj := range ls {
			ret = append(ret, obj)
		}
	}
	return ret, nil
}

func getAllDeviceModels(yamlDir string, schedulerNodePortIp string, imageRegistry string) ([]client.Object, error) {
	yamlList := []string{"ht-sensor-arm64", "rtmp-camera-arm64", "rtmp-camera", "vibration-sensor-arm64"}
	var ret []client.Object
	for _, ym := range yamlList {
		ls, err := deserializeAndRenderFromFile(filepath.Join(yamlDir, ym+".yaml"), "", "", "")
		if err != nil {
			return ret, err
		}
		for _, obj := range ls {
			ret = append(ret, obj)
		}
	}
	return ret, nil
}
