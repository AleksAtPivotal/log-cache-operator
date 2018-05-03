package logcache

import (
	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/sirupsen/logrus"
)

// Cleanup cleans up objects that Operator created
func Cleanup(lc *api.LogCache) error {
	logrus.Infof("Deleting Log Cache Deployment")

	lc = lc.DeepCopy()
	err := deleteLogCache(lc)
	if err != nil {
		return nil
	}

	err = deleteLogCacheNozzle(lc)
	if err != nil {
		return nil
	}

	err = deleteLogCacheScheduler(lc)
	if err != nil {
		return nil
	}

	err = DeletelogcacheConfigMap(lc)
	if err != nil {
		return nil
	}

	err = deleteLogCacheHeadlessService(lc)
	if err != nil {
		return nil
	}

	return nil
}
