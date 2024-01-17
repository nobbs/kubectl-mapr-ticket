package volume

import (
	apiSecret "github.com/nobbs/kubectl-mapr-ticket/pkg/api/secret"

	coreV1 "k8s.io/api/core/v1"
)

type Volume struct {
	Volume *coreV1.PersistentVolume
	Ticket *apiSecret.TicketSecret
}

// NewVolume creates a new Volume
func NewVolume(pv *coreV1.PersistentVolume) *Volume {
	return &Volume{
		Volume: pv,
	}
}

// Name returns the name of the volume
func (v *Volume) Name() string {
	if v.Volume == nil {
		return ""
	}

	return v.Volume.Name
}

// ClaimName returns the name of the PVC that is bound to the volume
func (v *Volume) ClaimName() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.ClaimRef == nil {
		return ""
	}

	return v.Volume.Spec.ClaimRef.Name
}

// ClaimNamespace returns the namespace of the PVC that is bound to the volume
func (v *Volume) ClaimNamespace() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.ClaimRef == nil {
		return ""
	}

	return v.Volume.Spec.ClaimRef.Namespace
}

// ClaimUID returns the volume path of the volume
func (v *Volume) VolumePath() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.CSI == nil {
		return ""
	}

	if v.Volume.Spec.CSI.VolumeAttributes == nil {
		return ""
	}

	value, ok := v.Volume.Spec.CSI.VolumeAttributes["volumePath"]
	if !ok {
		return ""
	}

	return value
}

// VolumeHandle returns the volume handle of the volume
func (v *Volume) VolumeHandle() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.CSI == nil {
		return ""
	}

	return v.Volume.Spec.CSI.VolumeHandle
}

// SecretName returns the name of the NodePublishSecretRef of the volume
func (v *Volume) SecretName() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.CSI == nil {
		return ""
	}

	if v.Volume.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.Volume.Spec.CSI.NodePublishSecretRef.Name
}

// SecretNamespace returns the namespace of the NodePublishSecretRef of the volume
func (v *Volume) SecretNamespace() string {
	if v.Volume == nil {
		return ""
	}

	if v.Volume.Spec.CSI == nil {
		return ""
	}

	if v.Volume.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.Volume.Spec.CSI.NodePublishSecretRef.Namespace
}
