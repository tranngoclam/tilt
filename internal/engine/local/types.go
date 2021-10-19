package local

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
)

type Cmd = v1alpha1.Cmd
type CmdList = v1alpha1.CmdList
type CmdStatus = v1alpha1.CmdStatus
type CmdSpec = v1alpha1.CmdSpec
type CmdStateWaiting = v1alpha1.CmdStateWaiting
type CmdStateTerminated = v1alpha1.CmdStateTerminated
type CmdStateRunning = v1alpha1.CmdStateRunning
type CmdServer = v1alpha1.CmdServer
type CmdServerSpec = v1alpha1.CmdServerSpec
type CmdServerStatus = v1alpha1.CmdServerStatus
type ObjectMeta = metav1.ObjectMeta
