package logcache

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/operator-framework/operator-sdk/pkg/sdk/query"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func logsupdateLogCacheStatus(l *api.LogCache, status *api.LogCacheStatus) error {
	// don't update the status if there aren't any changes.
	if reflect.DeepEqual(l.Status, *status) {
		return nil
	}
	l.Status = *status
	return action.Update(l)
}

// getLogCacheStatus retrieves the status of the log cache cluster for the given Custom Resource "lc",
// and it only succeeds if all of the nodes from log cache cluster are reachable.
func getLogCacheStatus(lc *api.LogCache) (*api.LogCacheStatus, error) {
	pods := &v1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}
	sel := labelsForLogCache(lc.Name)
	opt := &metav1.ListOptions{LabelSelector: labels.SelectorFromSet(sel).String()}
	err := query.List(lc.GetNamespace(), pods, query.WithListOptions(opt))
	if err != nil {
		return nil, fmt.Errorf("failed to get logcache's pods: %v", err)
	}

	var (
		initialized bool
		active      string
	)

	for _, p := range pods.Items {
		// if a pod is terminating, then we can't access the corresponding log cache node's status.
		// so we break away from here and return an error.
		if p.Status.Phase != v1.PodRunning || p.DeletionTimestamp != nil {
			return nil, errors.New("Log Cache pod is terminating")
		}
		initialized = true

	}

	return &api.LogCacheStatus{
		Phase:       api.ClusterPhaseRunning,
		Initialized: initialized,
		ServiceName: lc.GetName(),
		State:       active,
	}, nil
}
