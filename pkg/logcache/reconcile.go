package logcache

import (
	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/sirupsen/logrus"
)

// Reconcile reconciles the LogCache cluster's state to the spec specified by cr
// by preparing the TLS secrets, deploying the LogCache cluster,
// and finally updating the LogCache deployment if needed.
func Reconcile(lc *api.LogCache) (err error) {
	lc = lc.DeepCopy()

	// Simulate initializer.
	changed := lc.SetDefaults()
	if changed {
		logrus.Infof("Setting the defaults for the new Log Cache Cluster: %s", lc.Name)
		return action.Update(lc)
	}

	// After first time reconcile, phase will switch to "Running".
	if lc.Status.Phase == api.ClusterPhaseInitial {
		logrus.Infof("Setting up the initial cluster")
	}

	err = CreatelogcacheConfigMap(lc)
	if err != nil {
		return err
	}

	err = deployLogCacheNozzle(lc)
	if err != nil {
		return err
	}

	err = deployLogCacheScheduler(lc)
	if err != nil {
		return err
	}

	err = deployLogCache(lc)
	if err != nil {
		return err
	}

	ls, err := getLogCacheStatus(lc)
	if err != nil {
		return err
	}

	return logsupdateLogCacheStatus(lc, ls)
}
