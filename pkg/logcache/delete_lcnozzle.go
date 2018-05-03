package logcache

import (
	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteLogCacheNozzle(l *api.LogCache) error {
	logrus.Infof("Deleting Log Cache Nozzle Deployment")
	selector := labelsForLogCache(l.GetName())

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName() + "-nozzle",
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
	}

	err := action.Delete(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Infof("Error while deleting Log Cache Nozzle Deployment: %s", err)
		return err
	}
	return nil
}
