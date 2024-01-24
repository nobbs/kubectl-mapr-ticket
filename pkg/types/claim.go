package types

import (
	coreV1 "k8s.io/api/core/v1"
)

type PersistentVolumeClaim coreV1.PersistentVolumeClaim

type VolumeClaim struct {
	Claim  *PersistentVolumeClaim
	Volume *PersistentVolume
	Ticket *TicketSecret
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
	return c.Status.Phase == coreV1.ClaimBound
}
