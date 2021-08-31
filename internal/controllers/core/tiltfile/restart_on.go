package tiltfile

import (
	"context"
	"time"

	"github.com/tilt-dev/tilt/internal/sliceutils"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

// TODO(nick): This code is needed by anything with a RestartOnSpec.
// We should find a way to consolidate this.

// Fetch all the buttons that this object depends on.
func (r *Reconciler) buttons(ctx context.Context, obj *v1alpha1.Tiltfile) (map[string]*v1alpha1.UIButton, error) {
	buttonNames := []string{}

	restartOn := obj.Spec.RestartOn
	if restartOn != nil {
		buttonNames = append(buttonNames, restartOn.UIButtons...)
	}

	result := make(map[string]*v1alpha1.UIButton, len(buttonNames))
	for _, n := range buttonNames {
		_, exists := result[n]
		if exists {
			continue
		}

		b := &v1alpha1.UIButton{}
		err := r.ctrlClient.Get(ctx, types.NamespacedName{Name: n}, b)
		if err != nil {
			return nil, err
		}
		result[n] = b
	}
	return result, nil
}

// Fetch all the filewatches that this object depends on.
func (r *Reconciler) fileWatches(ctx context.Context, obj *v1alpha1.Tiltfile) (map[string]*v1alpha1.FileWatch, error) {
	restartOn := obj.Spec.RestartOn
	if restartOn == nil {
		return nil, nil
	}

	result := make(map[string]*v1alpha1.FileWatch, len(restartOn.FileWatches))
	for _, n := range restartOn.FileWatches {
		fw := &v1alpha1.FileWatch{}
		err := r.ctrlClient.Get(ctx, types.NamespacedName{Name: n}, fw)
		if err != nil {
			return nil, err
		}
		result[n] = fw
	}
	return result, nil
}

// Fetch the last time a restart was requested from this target's dependencies.
func (r *Reconciler) lastRestartEvent(restartOn *v1alpha1.RestartOnSpec, fileWatches map[string]*v1alpha1.FileWatch, buttons map[string]*v1alpha1.UIButton) time.Time {
	cur := time.Time{}
	if restartOn == nil {
		return cur
	}

	for _, fwn := range restartOn.FileWatches {
		fw, ok := fileWatches[fwn]
		if !ok {
			// ignore missing filewatches
			continue
		}
		lastEventTime := fw.Status.LastEventTime
		if lastEventTime.Time.After(cur) {
			cur = lastEventTime.Time
		}
	}

	for _, bn := range restartOn.UIButtons {
		b, ok := buttons[bn]
		if !ok {
			// ignore missing buttons
			continue
		}
		lastEventTime := b.Status.LastClickedAt
		if lastEventTime.Time.After(cur) {
			cur = lastEventTime.Time
		}
	}

	return cur
}

// Fetch the set of files that have changed since the given timestamp.
// We err on the side of undercounting (i.e., skipping files that may have triggered
// this build but are not sure).
func (r *Reconciler) filesChanged(restartOn *v1alpha1.RestartOnSpec, fileWatches map[string]*v1alpha1.FileWatch, lastBuild time.Time) []string {
	filesChanged := []string{}
	for _, fwn := range restartOn.FileWatches {
		fw, ok := fileWatches[fwn]
		if !ok {
			// ignore missing filewatches
			continue
		}

		// Add files so that the most recent files are first.
		for i := len(fw.Status.FileEvents) - 1; i >= 0; i-- {
			e := fw.Status.FileEvents[i]
			if e.Time.Time.After(lastBuild) {
				filesChanged = append(filesChanged, e.SeenFiles...)
			}
		}
	}
	return sliceutils.DedupedAndSorted(filesChanged)
}
