package logcache

import (
	"fmt"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatelogcacheConfigMap Creates a configmap to store configuration
func CreatelogcacheConfigMap(l *api.LogCache) error {
	m := make(map[string]string)
	m["NODE_ADDRS"] = GetLogCacheNodeAddresses(l)

	configmap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName() + "-config",
			Namespace: l.GetNamespace(),
		},
		Data: m,
	}

	err := action.Create(configmap)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Error while creating Log Cache Configmap: %v", err)
	}
	if err != nil && apierrors.IsAlreadyExists(err) {

		update := action.Update(configmap)
		if update != nil {
			return fmt.Errorf("Error while Updating Log Cache Configmap: %v", update)
		}
	}

	return nil
}

// DeletelogcacheConfigMap Deletes a configmap to store configuration
func DeletelogcacheConfigMap(l *api.LogCache) error {
	configmap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName() + "-config",
			Namespace: l.GetNamespace(),
		},
	}

	err := action.Delete(configmap)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Error while deleting Log Cache Configmap: %v", err)
	}

	return nil
}

// GetLogCacheNodeAddresses returns a string output of LogCache cluster service addresses
func GetLogCacheNodeAddresses(l *api.LogCache) (NodeAddress string) {
	nodecount := l.Spec.LogCachePod.Nodes
	lcname := l.GetName()
	lcns := l.GetNamespace()
	var i int32
	NodeAddress = ""

	for i = 0; i < nodecount; i++ {
		if NodeAddress == "" {
			NodeAddress = lcname + "-" + fmt.Sprintf("%v", i) + "." + lcname + "." + lcns + "." + clusterservicedomain + fmt.Sprintf(":%v", logcacheClientPort)
		} else {
			NodeAddress = NodeAddress + "," + lcname + "-" + fmt.Sprintf("%v", i) + "." + lcname + "." + lcns + "." + clusterservicedomain + fmt.Sprintf(":%v", logcacheClientPort)
		}

	}
	return NodeAddress
}
