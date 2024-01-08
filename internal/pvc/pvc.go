package pvc

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/volume"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Lister struct {
	client    kubernetes.Interface
	namespace string
}

type ListerOption func(*Lister)

func NewLister(client kubernetes.Interface, namespace string, opts ...ListerOption) *Lister {
	l := &Lister{
		client:    client,
		namespace: namespace,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Lister) Run() ([]coreV1.PersistentVolumeClaim, error) {
	// Get all PVCs in the namespace
	pvcs, err := l.client.CoreV1().PersistentVolumeClaims(l.namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter the PVCs to only bound ones
	filtered := l.FilterBoundPVCs(pvcs.Items)

	// Filter the PVCs to only ones that are bound to MapR CSI-based PVs
	filtered = l.FilterPVCsToMaprCSI(filtered)

	return filtered, nil
}

func (l *Lister) FilterBoundPVCs(pvcs []coreV1.PersistentVolumeClaim) []coreV1.PersistentVolumeClaim {
	filtered := make([]coreV1.PersistentVolumeClaim, 0, len(pvcs))
	for _, pvc := range pvcs {
		if pvc.Status.Phase == coreV1.ClaimBound {
			filtered = append(filtered, pvc)
		}
	}

	return filtered
}

func (l *Lister) FilterPVCsToMaprCSI(pvcs []coreV1.PersistentVolumeClaim) []coreV1.PersistentVolumeClaim {
	filtered := make([]coreV1.PersistentVolumeClaim, 0, len(pvcs))
	for _, pvc := range pvcs {
		if l.isPVCBoundToMaprPV(pvc) {
			filtered = append(filtered, pvc)
		}
	}

	return filtered
}

func (l *Lister) isPVCBoundToMaprPV(pvc coreV1.PersistentVolumeClaim) bool {
	// Get the PV that the PVC is bound to
	pv, err := l.client.CoreV1().PersistentVolumes().Get(context.TODO(), pvc.Spec.VolumeName, metaV1.GetOptions{})
	if err != nil {
		return false
	}

	// Check if the PV is a MapR CSI PV
	return volume.IsMaprCSIBased(pv)
}
