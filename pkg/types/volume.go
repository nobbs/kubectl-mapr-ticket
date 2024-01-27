// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package types

import (
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
)

const (
	MaprCSIProvisionerKDF    = "com.mapr.csi-kdf"
	MaprCSIProvisionerNFSKDF = "com.mapr.csi-nfskdf"
)

var (
	// MaprCSIProvisioners is a list of the default MapR CSI provisioners
	// that we support.
	MaprCSIProvisioners = []string{
		MaprCSIProvisionerKDF,
		MaprCSIProvisionerNFSKDF,
	}
)

// PersistentVolume is a wrapper around coreV1.PersistentVolume that provides additional
// functionality.
type PersistentVolume coreV1.PersistentVolume

// MaprVolume is a wrapper around a PersistentVolume that provides additional functionality.
type MaprVolume struct {
	Volume *PersistentVolume
	Ticket *MaprSecret
}

// GetName returns the name of the volume
func (v *PersistentVolume) GetName() string {
	if v == nil {
		return ""
	}

	return v.Name
}

// GetClaimName returns the name of the PVC that is bound to the volume
func (v *PersistentVolume) GetClaimName() string {
	if v == nil || v.Spec.ClaimRef == nil {
		return ""
	}

	return v.Spec.ClaimRef.Name
}

// GetClaimNamespace returns the namespace of the PVC that is bound to the volume
func (v *PersistentVolume) GetClaimNamespace() string {
	if v == nil || v.Spec.ClaimRef == nil {
		return ""
	}

	return v.Spec.ClaimRef.Namespace
}

// ClaimUID returns the volume path of the volume
func (v *PersistentVolume) GetVolumePath() string {
	if v == nil || v.Spec.CSI == nil || v.Spec.CSI.VolumeAttributes == nil {
		return ""
	}

	value, ok := v.Spec.CSI.VolumeAttributes["volumePath"]
	if !ok {
		return ""
	}

	return value
}

// GetVolumeHandle returns the volume handle of the volume
func (v *PersistentVolume) GetVolumeHandle() string {
	if v == nil || v.Spec.CSI == nil {
		return ""
	}

	return v.Spec.CSI.VolumeHandle
}

// GetSecretName returns the name of the NodePublishSecretRef of the volume
func (v *PersistentVolume) GetSecretName() string {
	if v == nil || v.Spec.CSI == nil {
		return ""
	}

	if v.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.Spec.CSI.NodePublishSecretRef.Name
}

// GetSecretNamespace returns the namespace of the NodePublishSecretRef of the volume
func (v *PersistentVolume) GetSecretNamespace() string {
	if v == nil || v.Spec.CSI == nil || v.Spec.CSI.NodePublishSecretRef == nil {
		return ""
	}

	return v.Spec.CSI.NodePublishSecretRef.Namespace
}

// IsMaprCSIBased returns true if the volume is provisioned by one of the MapR CSI provisioners and
// false otherwise.
func (v *PersistentVolume) IsMaprCSIBased() bool {
	if v == nil || v.Spec.CSI == nil {
		return false
	}

	// Check if the volume is provisioned by one of the MapR CSI provisioners
	for _, provisioner := range MaprCSIProvisioners {
		if v.Spec.CSI.Driver == provisioner {
			return true
		}
	}

	return false
}

// UsesSecret returns true if the volume uses the specified secret and false otherwise. If the
// secret name is equal to the value of SecretAll, all volumes that use a secret from the specified
// namespace are returned. If the secret namespace is equal to the value of NamespaceAll, basically
// all volumes that any secret in any namespace will evaluate to true.
func (volume *PersistentVolume) UsesSecret(namespace, name string) bool {
	if volume == nil || volume.Spec.CSI == nil || volume.Spec.CSI.NodePublishSecretRef == nil {
		return false
	}

	// Check if we want secrets from all namespaces
	if namespace == util.NamespaceAll {
		return true
	}

	// Check if we want all secrets from the specified namespace
	if name == util.SecretAll && volume.Spec.CSI.NodePublishSecretRef.Namespace == namespace {
		return true
	}

	// Check if the volume uses the specified secret
	if volume.Spec.CSI.NodePublishSecretRef.Name != name {
		return false
	}

	// Check if the volume uses the specified namespace
	if volume.Spec.CSI.NodePublishSecretRef.Namespace != namespace {
		return false
	}

	return true
}
