package logcache

import (
	"fmt"
	"path/filepath"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deployLogCache deploys a Log Cache service.
// deployLogCache is a multi-steps process. It creates the deployment, the service and
// other related Kubernetes objects for Log Cache. Any intermediate step can fail.
//
// deployLogCache is idempotent. If an object already exists, this function will ignore creating
// it and return no error. It is safe to retry on this function.
func deployLogCache(l *api.LogCache) error {
	selector := labelsForLogCache(l.GetName())
	podTempl := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName(),
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{logcacheContainer(l), logcachegatewayContainer(l)},
			Volumes: []v1.Volume{{
				Name: logcacheCertVolumeName,
				VolumeSource: v1.VolumeSource{
					Projected: &v1.ProjectedVolumeSource{
						Sources: []v1.VolumeProjection{{
							Secret: &v1.SecretProjection{
								LocalObjectReference: v1.LocalObjectReference{
									// [todo] -- this is hardcoded for now
									Name: logcacheTLSSecret,
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

	d := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.GetName(),
			Namespace: l.GetNamespace(),
			Labels:    selector,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &l.Spec.LogCachePod.Nodes,
			Selector:    &metav1.LabelSelector{MatchLabels: selector},
			Template:    podTempl,
			ServiceName: l.GetName(),
		},
	}
	err := action.Create(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Error while creating Log Cache Deployment: %v", err)
	}

	// If Object already exists check to if it's a scale-out / in event
	if err != nil && apierrors.IsAlreadyExists(err) {
		currentd := &appsv1.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StatefulSet",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      l.GetName(),
				Namespace: l.GetNamespace(),
				Labels:    selector,
			}}
		sync := query.Get(currentd)
		if sync != nil {
			return fmt.Errorf("failed to get deployment: %v", sync)
		}
		nodes := l.Spec.LogCachePod.Nodes

		if *currentd.Spec.Replicas != nodes {
			logrus.Infof("Destroying and recreating the Log Cache Statefulset")
			action.Delete(d)
			action.Create(d)
			return nil
		}

		update := action.Update(d)
		if update != nil {
			return fmt.Errorf("Error while Updating Log Cache StateFulset: %v", update)
		}
	}

	err = deployLogCacheHeadlessService(l)
	if err != nil {
		return fmt.Errorf("Error while creating LogCache Service : %v", err)
	}

	return nil
}

func logcacheContainer(l *api.LogCache) v1.Container {
	return v1.Container{
		Name:  "log-cache",
		Image: fmt.Sprintf("%s:%s", l.Spec.LogCachePod.BaseImage, l.Spec.LogCachePod.Version),

		Ports: []v1.ContainerPort{{
			Name:          logcacheClientPortName,
			ContainerPort: int32(logcacheClientPort),
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
		}},
	}
}

func logcachegatewayContainer(l *api.LogCache) v1.Container {
	return v1.Container{
		Name:  "log-cache-gateway",
		Image: fmt.Sprintf("%s:%s", l.Spec.LogCachePod.GatewayBaseImage, l.Spec.LogCachePod.GatewayVersion),
		Ports: []v1.ContainerPort{{
			Name:          logcacheGatewayPortName,
			ContainerPort: int32(logcacheGatewayPort),
		}},
		VolumeMounts: []v1.VolumeMount{{
			Name:      logcacheCertVolumeName,
			MountPath: filepath.Dir(logcacheCertPath),
		}},
		Env: []v1.EnvVar{
			{
				Name:  "LOG_CACHE_ADDR",
				Value: fmt.Sprintf("localhost:%v", logcacheClientPort),
			},
			{
				Name:  "GROUP_READER_ADDR",
				Value: fmt.Sprintf("localhost:%v", logcacheClientPort),
			}},
	}
}
