package logcache

import (
	"fmt"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deployLogCacheHeadlessService Deploys headless service for LogCache Statefulset
func deployLogCacheHeadlessService(l *api.LogCache) error {
	selector := labelsForLogCache(l.GetName())

	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName(),
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Ports: []v1.ServicePort{
				{
					Name: logcacheClientPortName,
					Port: logcacheClientPort,
				},
				{
					Name: logcacheGatewayPortName,
					Port: logcacheGatewayPort,
				},
			},
			ClusterIP: "None",
			Type:      "ClusterIP",
		},
	}

	err := action.Create(svc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create logcache service: %v", err)
	}
	return nil
}

// deleteLogCacheHeadlessService Deletes headless service for LogCache Statefulset
func deleteLogCacheHeadlessService(l *api.LogCache) error {
	selector := labelsForLogCache(l.GetName())

	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName(),
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
	}

	err := action.Delete(svc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create logcache service: %v", err)
	}
	return nil
}
