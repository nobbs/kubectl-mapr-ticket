package types

import (
	coreV1 "k8s.io/api/core/v1"
)

// PersistentVolumeClaim is a wrapper around coreV1.PersistentVolumeClaim that provides additional
// functionality.
type PersistentVolumeClaim coreV1.PersistentVolumeClaim

// MaprVolumeClaim is used to store a claim, the volume it is bound to, and the ticket used by that
// volume
type MaprVolumeClaim struct {
	Claim  *PersistentVolumeClaim
	Volume *PersistentVolume
	Ticket *MaprSecret
}

// GetNamespace returns the namespace of the claim
func (c *PersistentVolumeClaim) GetNamespace() string {
	if c == nil {
		return ""
	}

	return c.Namespace
}

// GetName returns the name of the claim
func (c *PersistentVolumeClaim) GetName() string {
	if c == nil {
		return ""
	}

	return c.Name
}

// IsBound returns true if the claim is bound to a volume
func (c *PersistentVolumeClaim) IsBound() bool {
	if c == nil {
		return false
	}

	return c.Status.Phase == coreV1.ClaimBound
}
