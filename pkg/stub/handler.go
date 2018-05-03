package stub

import (
	"github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/alekssaul/logcache-operator/pkg/logcache"
	"github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	"github.com/operator-framework/operator-sdk/pkg/sdk/types"
)

func NewHandler() handler.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.LogCache:
		if event.Deleted {
			logrus.Infof("Logcache CRD is removed, cleaning up")
			err := logcache.Cleanup(o)
			if err != nil {
				logrus.Infof("Error: %s", err)
				return err
			}

			return nil
		}

		err := logcache.Reconcile(o)
		if err != nil {
			logrus.Infof("Error: %s", err)
		}
		return err
	}
	return nil
}
