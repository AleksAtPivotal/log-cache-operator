package logcache

import (
	"fmt"
	"path/filepath"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// deployLogCacheScheduler deploys a Log Cache Scheduler.
func deployLogCacheScheduler(l *api.LogCache) error {
	selector := labelsForLogCache(l.GetName())
	podTempl := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName() + "-scheduler",
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{logcacheSchedulerContainer(l)},
			Volumes: []v1.Volume{{
				Name: logcacheCertVolumeName,
				VolumeSource: v1.VolumeSource{
					Projected: &v1.ProjectedVolumeSource{
						Sources: []v1.VolumeProjection{{
							Secret: &v1.SecretProjection{
								LocalObjectReference: v1.LocalObjectReference{
									// [todo] -- this is hardcoded for now
									Name: logcacheSchedulerTLSSecret,
								},
							},
						}},
					},
				},
			}},
			SecurityContext: &v1.PodSecurityContext{
				RunAsUser:    func(i int64) *int64 { return &i }(9000),
				RunAsNonRoot: func(b bool) *bool { return &b }(true),
				FSGroup:      func(i int64) *int64 { return &i }(9000),
			},
		},
	}

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName() + "-scheduler",
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &l.Spec.LogCacheScheduler.Nodes,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: podTempl,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
					MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				},
			},
		},
	}

	err := action.Create(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Error while creating Log Cache Scheduler Deployment: %v", err)
	}
	if err != nil && apierrors.IsAlreadyExists(err) {
		update := action.Update(d)
		if update != nil {
			return fmt.Errorf("Error while Updating Log Cache Scheduler Deployment: %v", update)
		}
	}
	return nil
}

func logcacheSchedulerContainer(l *api.LogCache) v1.Container {

	return v1.Container{
		Name:  "scheduler",
		Image: fmt.Sprintf("%s:%s", l.Spec.LogCacheScheduler.BaseImage, l.Spec.LogCacheScheduler.Version),
		Ports: []v1.ContainerPort{{
			Name:          logcacheSchedulerPortName,
			ContainerPort: int32(logcacheSchedulerPort),
		}},
		VolumeMounts: []v1.VolumeMount{{
			Name:      logcacheCertVolumeName,
			MountPath: filepath.Dir(logcacheCertPath),
		}},
		Env: []v1.EnvVar{{
			Name: "NODE_ADDRS",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					Key: "NODE_ADDRS",
					LocalObjectReference: v1.LocalObjectReference{
						Name: l.GetName() + "-config",
					},
				},
			},
		}, {
			Name:  "GROUP_READER_NODE_ADDRS",
			Value: l.GetName() + fmt.Sprintf(":%v", logcacheClientPort),
		}},
	}
}
