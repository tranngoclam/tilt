package k8sconv

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tilt-dev/tilt/internal/k8s"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// A KubernetesResource exposes a high-level status that summarizes
// the Pods we care about in a KubernetesDiscovery.
//
// If we have a KubernetesApply, KubernetesResource will use that
// to narrow down the list of pods to only the pods we care about
// for the current Apply.
//
// KubernetesResource is intended to be a non-stateful object (i.e., it is
// immutable and its status can be inferred from the state of child
// objects.)
//
// Long-term, this may become an explicit API server object, but
// for now it's intended to provide an API-server compatible
// layer around KubernetesDiscovery + KubernetesApply.
type KubernetesResource struct {
	Discovery   *v1alpha1.KubernetesDiscovery
	ApplyStatus *v1alpha1.KubernetesApplyStatus

	// A set of properties we use to determine which pods in Discovery
	// belong to the current Apply.
	ApplyFilter *KubernetesApplyFilter
}

func NewKubernetesResource(discovery *v1alpha1.KubernetesDiscovery, status *v1alpha1.KubernetesApplyStatus) (*KubernetesResource, error) {
	var filter *KubernetesApplyFilter
	var err error
	if status != nil {
		filter, err = NewKubernetesApplyFilter(status)
		if err != nil {
			return nil, err
		}
	}

	return &KubernetesResource{Discovery: discovery, ApplyStatus: status, ApplyFilter: filter}, nil

}

type KubernetesApplyFilter struct {
	// DeployedRefs are references to the objects that we deployed to a Kubernetes cluster.
	DeployedRefs []v1.ObjectReference

	// Hashes of the pod template specs that we deployed to a Kubernetes cluster.
	PodTemplateSpecHashes []k8s.PodTemplateSpecHash
}

func NewKubernetesApplyFilter(status *v1alpha1.KubernetesApplyStatus) (*KubernetesApplyFilter, error) {
	deployed, err := k8s.ParseYAMLFromString(status.ResultYAML)
	if err != nil {
		return nil, err
	}

	podTemplateSpecHashes := []k8s.PodTemplateSpecHash{}
	for _, entity := range deployed {
		if entity.UID() == "" {
			return nil, fmt.Errorf("Entity not deployed correctly: %v", entity)
		}
		hs, err := k8s.ReadPodTemplateSpecHashes(entity)
		if err != nil {
			return nil, errors.Wrap(err, "reading pod template spec hashes")
		}
		podTemplateSpecHashes = append(podTemplateSpecHashes, hs...)
	}
	return &KubernetesApplyFilter{
		DeployedRefs:          k8s.ToRefList(deployed),
		PodTemplateSpecHashes: podTemplateSpecHashes,
	}, nil
}
