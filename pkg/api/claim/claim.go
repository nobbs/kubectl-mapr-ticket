package claim

import (
	apiSecret "github.com/nobbs/kubectl-mapr-ticket/pkg/api/secret"

	coreV1 "k8s.io/api/core/v1"
)

type VolumeClaim struct {
	Claim  *coreV1.PersistentVolumeClaim
	Volume *coreV1.PersistentVolume
	Ticket *apiSecret.TicketSecret
}

// NewVolumeClaim returns a new VolumeClaim from the given PVC.
func NewVolumeClaim(pvc *coreV1.PersistentVolumeClaim) *VolumeClaim {
	return &VolumeClaim{
		Claim: pvc,
	}
}

// IsClaimBound returns true if the claim is bound to a volume
func (vc *VolumeClaim) IsClaimBound() bool {
	return vc.Claim.Status.Phase == coreV1.ClaimBound
}
