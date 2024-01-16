package util

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"

	coreV1 "k8s.io/api/core/v1"
)

type TicketSecret struct {
	Secret *coreV1.Secret `json:"originalSecret"`
	Ticket *ticket.Ticket `json:"parsedTicket"`
	NumPVC uint32         `json:"numPVC"`
}

type VolumeClaim struct {
	PVC *coreV1.PersistentVolumeClaim
	PV  *coreV1.PersistentVolume
}

type Volume struct {
	PV     *coreV1.PersistentVolume
	Ticket *TicketSecret
}

// Name returns the name of the volume
func (v *Volume) Name() string {
	if v.PV == nil {
		return ""
	}

	return v.PV.Name
}

// ClaimName returns the name of the PVC that is bound to the volume
func (v *Volume) ClaimName() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.ClaimRef == nil {
		return ""
	}

	return v.PV.Spec.ClaimRef.Name
}

// ClaimNamespace returns the namespace of the PVC that is bound to the volume
func (v *Volume) ClaimNamespace() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.ClaimRef == nil {
		return ""
	}

	return v.PV.Spec.ClaimRef.Namespace
}

// ClaimUID returns the volume path of the volume
func (v *Volume) VolumePath() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.CSI == nil {
		return ""
	}

	if v.PV.Spec.CSI.VolumeAttributes == nil {
		return ""
	}

	value, ok := v.PV.Spec.CSI.VolumeAttributes["volumePath"]
	if !ok {
		return ""
	}

	return value
}

// VolumeHandle returns the volume handle of the volume
func (v *Volume) VolumeHandle() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.CSI == nil {
		return ""
	}

	return v.PV.Spec.CSI.VolumeHandle
}

// SecretName returns the name of the NodePublishSecretRef of the volume
func (v *Volume) SecretName() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.CSI == nil {
		return ""
	}

	if v.PV.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.PV.Spec.CSI.NodePublishSecretRef.Name
}

// SecretNamespace returns the namespace of the NodePublishSecretRef of the volume
func (v *Volume) SecretNamespace() string {
	if v.PV == nil {
		return ""
	}

	if v.PV.Spec.CSI == nil {
		return ""
	}

	if v.PV.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.PV.Spec.CSI.NodePublishSecretRef.Namespace
}
