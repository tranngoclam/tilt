package tiltfile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tilt-dev/tilt/internal/store/tiltfiles"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	"github.com/tilt-dev/tilt/pkg/model"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// TODO(nick): This code is generally useful by anything that needs to read
// the trigger queue.
func (r *Reconciler) triggerQueue(ctx context.Context) (*v1alpha1.ConfigMap, error) {
	var cm v1alpha1.ConfigMap
	err := r.ctrlClient.Get(ctx, types.NamespacedName{Name: tiltfiles.TriggerQueueConfigMapName}, &cm)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	return &cm, nil
}

func (r *Reconciler) inTriggerQueue(cm *v1alpha1.ConfigMap, nn types.NamespacedName) bool {
	name := nn.Name
	for k, v := range cm.Data {
		if strings.HasSuffix(k, "-reason-code") {
			continue
		}
		if v == name {
			return true
		}
	}
	return false
}

func (r *Reconciler) triggerQueueReason(cm *v1alpha1.ConfigMap, nn types.NamespacedName) model.BuildReason {
	name := nn.Name
	for k, v := range cm.Data {
		if strings.HasSuffix(k, "-reason-code") {
			continue
		}

		if v != name {
			continue
		}

		reasonCode := cm.Data[fmt.Sprintf("%s-reason-code", k)]
		i, err := strconv.Atoi(reasonCode)
		if err != nil {
			return model.BuildReasonFlagTriggerUnknown
		}
		return model.BuildReason(i)
	}
	return model.BuildReasonNone
}
