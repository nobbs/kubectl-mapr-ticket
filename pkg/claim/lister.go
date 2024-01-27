// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

// Package claim implements the persistent volume claim lister. It is responsible for listing all
// persistent volume claims in the cluster that are refering to MapR-backed persistent volumes.
//
// The lister is implemented as a chain of filters and collectors. The filters are used to filter
// out volume claims, ie. to only keep those that are bound and backed by a MapR CSI provisioner.
// The collectors are used to collect additional information about the volume claims, ie. to collect
// the PV for each PVC and the MapR ticket for each PV. This data is then used to print the volume
// claims in a human-readable tabular format.
package claim

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type secretLister interface {
	List() ([]types.MaprSecret, error)
}

// Lister is the struct that is used to list volume claims refering to MapR-backed persistent
// volumes in the cluster.
type Lister struct {
	client    kubernetes.Interface
	namespace string

	secretLister secretLister
	sortBy       []SortOption

	volumeClaims []types.MaprVolumeClaim
}

// NewLister creates a new volume claim lister. It requires a Kubernetes client and a namespace
// to operate on. It also accepts a list of options that can be used to configure the lister.
func NewLister(client kubernetes.Interface, namespace string, opts ...ListerOption) *Lister {
	l := &Lister{
		client:    client,
		namespace: namespace,
		sortBy:    DefaultSortBy,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// List returns a list of volume claims that are provisioned by one of the MapR CSI provisioners.
func (l *Lister) List() ([]types.MaprVolumeClaim, error) {
	if err := l.getClaims(); err != nil {
		return nil, err
	}

	l.filterClaimsBoundOnly().
		collectVolumes().
		filterClaimsMaprCSI().
		collectTickets().
		sort()

	return l.volumeClaims, nil
}

// getClaims returns a list of all PVCs in the cluster
func (l *Lister) getClaims() error {
	claims, err := l.client.CoreV1().PersistentVolumeClaims(l.namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	volumeClaims := make([]types.MaprVolumeClaim, 0, len(claims.Items))

	for i := range claims.Items {
		volumeClaims = append(
			volumeClaims,
			types.MaprVolumeClaim{
				Claim: (*types.PersistentVolumeClaim)(&claims.Items[i]),
			},
		)
	}

	l.volumeClaims = volumeClaims

	return nil
}

// filterClaimsBoundOnly filters PVCs to those that are bound.
func (l *Lister) filterClaimsBoundOnly() *Lister {
	filtered := make([]types.MaprVolumeClaim, 0, len(l.volumeClaims))

	for _, volumeClaim := range l.volumeClaims {
		if volumeClaim.Claim.IsBound() {
			filtered = append(filtered, volumeClaim)
		}
	}

	l.volumeClaims = filtered

	return l
}

// filterClaimsMaprCSI filters PVCs to those that are provisioned by one of the MapR CSI
// provisioners.
func (l *Lister) filterClaimsMaprCSI() *Lister {
	filtered := make([]types.MaprVolumeClaim, 0, len(l.volumeClaims))

	for _, volumeClaim := range l.volumeClaims {
		if volumeClaim.Volume.IsMaprCSIBased() {
			filtered = append(filtered, volumeClaim)
		}
	}

	l.volumeClaims = filtered

	return l
}

// collectVolumes collects the PV for each PVC.
func (l *Lister) collectVolumes() *Lister {
	// Get all PVs in the cluster
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return l
	}

	// Lookup the PV for each PVC
	lookupPV := func(pvc *types.PersistentVolumeClaim) (*coreV1.PersistentVolume, bool) {
		for _, pv := range pvs.Items {
			if pv.Spec.CSI != nil && pv.Name == pvc.Spec.VolumeName {
				return &pv, true
			}
		}

		return nil, false
	}

	filtered := make([]types.MaprVolumeClaim, 0, len(l.volumeClaims))
	for _, volumeClaim := range l.volumeClaims {
		if pv, ok := lookupPV(volumeClaim.Claim); ok {
			volumeClaim.Volume = (*types.PersistentVolume)(pv)
			filtered = append(filtered, volumeClaim)
		}
	}

	l.volumeClaims = filtered

	return l
}

// collectTickets collects the MapR tickets for each PVC, if available.
func (l *Lister) collectTickets() *Lister {
	// return early if there is no secret lister
	if l.secretLister == nil {
		return l
	}

	// return early if there are no volume claims
	if len(l.volumeClaims) == 0 {
		return l
	}

	// collect all tickets via the secret lister
	tickets, err := l.secretLister.List()
	if err != nil {
		return l
	}

	// lookup the ticket for each volume claim
	lookupTicket := func(volumeClaim *types.MaprVolumeClaim) *types.MaprSecret {
		for _, ticket := range tickets {
			if ticket.Secret.Namespace == volumeClaim.Volume.GetSecretNamespace() &&
				ticket.Secret.Name == volumeClaim.Volume.GetSecretName() {
				return &ticket
			}
		}

		return nil
	}

	for i := range l.volumeClaims {
		volumeClaim := &l.volumeClaims[i]
		volumeClaim.Ticket = lookupTicket(volumeClaim)
	}

	return l
}
