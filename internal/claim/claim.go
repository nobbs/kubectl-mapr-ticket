package claim

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/volume"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type VolumeClaim struct {
	pvc *coreV1.PersistentVolumeClaim
	pv  *coreV1.PersistentVolume
}

type Lister struct {
	client    kubernetes.Interface
	namespace string

	volumeClaims []VolumeClaim
}

func NewLister(client kubernetes.Interface, namespace string) *Lister {
	l := &Lister{
		client:    client,
		namespace: namespace,
	}

	return l
}

func (l *Lister) List() ([]VolumeClaim, error) {
	// Get all PVCs in the namespace
	if err := l.getClaims(); err != nil {
		return nil, err
	}

	// Filter the PVCs to only those that are bound and MapR CSI based
	l.filterClaimsBoundOnly().
		filterClaimsMaprCSI()

	return l.volumeClaims, nil
}

func (l *Lister) getClaims() error {
	claims, err := l.client.CoreV1().PersistentVolumeClaims(l.namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	volumeClaims := make([]VolumeClaim, 0, len(claims.Items))

	for i := range claims.Items {
		volumeClaims = append(volumeClaims, VolumeClaim{
			pvc: &claims.Items[i],
		})
	}

	l.volumeClaims = volumeClaims

	return nil
}

func (l *Lister) filterClaimsBoundOnly() *Lister {
	filtered := make([]VolumeClaim, 0, len(l.volumeClaims))

	for _, volumeClaim := range l.volumeClaims {
		if volumeClaim.pvc.Status.Phase == coreV1.ClaimBound {
			filtered = append(filtered, volumeClaim)
		}
	}

	l.volumeClaims = filtered

	return l
}

func (l *Lister) filterClaimsMaprCSI() *Lister {
	// Get all PVs in the cluster
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return l
	}

	lookupPV := func(pvc *coreV1.PersistentVolumeClaim) *coreV1.PersistentVolume {
		for _, pv := range pvs.Items {
			if pv.Spec.CSI != nil && pv.Name == pvc.Spec.VolumeName {
				return &pv
			}
		}

		return nil
	}

	filtered := make([]VolumeClaim, 0, len(l.volumeClaims))
	for _, volumeClaim := range l.volumeClaims {
		pv := lookupPV(volumeClaim.pvc)
		if pv == nil {
			continue
		}

		if !volume.IsMaprCSIBased(pv) {
			continue
		}

		volumeClaim.pv = pv

		filtered = append(filtered, volumeClaim)
	}

	l.volumeClaims = filtered

	return l
}
