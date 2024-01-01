package volume

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// maprCSIProvisioners is a list of the default MapR CSI provisioners
	// that we support.
	maprCSIProvisioners = []string{
		"com.mapr.csi-kdf",
		"com.mapr.csi-nfskdf",
	}
)

type Lister struct {
	client     kubernetes.Interface
	secretName string
	namespace  string

	allSecrets bool
}

type ListerOption func(*Lister)

func WithAllSecrets() ListerOption {
	return func(l *Lister) {
		l.allSecrets = true
	}
}

func NewLister(client kubernetes.Interface, secretName string, namespace string, opts ...ListerOption) *Lister {
	l := &Lister{
		client:     client,
		secretName: secretName,
		namespace:  namespace,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Lister) Run() ([]coreV1.PersistentVolume, error) {
	// Unfortunately, we have to list all persistent volumes and filter them
	// ourselves, because there is no way to filter them by label selector.
	volumes, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter the volumes to only MapR CSI-based ones
	filtered := FilterVolumesToMaprCSI(volumes.Items)

	// If we are listing volumes for all secrets in the namespace, let's
	// filter the volumes to only ones that use a NodePublishSecretRef in
	// the namespace. Otherwise, let's filter the volumes to only ones that
	// use the specified secret.
	if l.allSecrets {
		filtered = filterVolumeUsesTicketFromNamespace(filtered, l.namespace)
	} else {
		filtered = filterVolumeUsesTicket(filtered, l.secretName, l.namespace)
	}

	return filtered, nil
}

func FilterVolumesToMaprCSI(volumes []coreV1.PersistentVolume) []coreV1.PersistentVolume {
	var filtered []coreV1.PersistentVolume

	for _, volume := range volumes {
		if isMaprCSIBased(&volume) {
			filtered = append(filtered, volume)
		}
	}

	return filtered
}

func filterVolumeUsesTicketFromNamespace(volumes []coreV1.PersistentVolume, namespace string) []coreV1.PersistentVolume {
	var filtered []coreV1.PersistentVolume

	for _, volume := range volumes {
		if usesTicketFromNamespace(&volume, namespace) {
			filtered = append(filtered, volume)
		}
	}

	return filtered
}

func filterVolumeUsesTicket(volumes []coreV1.PersistentVolume, secretName, namespace string) []coreV1.PersistentVolume {
	var filtered []coreV1.PersistentVolume

	for _, volume := range volumes {
		if UsesTicket(&volume, secretName, namespace) {
			filtered = append(filtered, volume)
		}
	}

	return filtered
}

func usesTicketFromNamespace(volume *coreV1.PersistentVolume, namespace string) bool {
	// Check if the volume uses a CSI driver
	if volume.Spec.CSI == nil {
		return false
	}

	// Check if the volume uses a NodePublishSecretRef
	if volume.Spec.CSI.NodePublishSecretRef == nil {
		return false
	}

	// Check if the volume uses a NodePublishSecretRef in the specified namespace
	if volume.Spec.CSI.NodePublishSecretRef.Namespace != namespace {
		return false
	}

	return true
}

func UsesTicket(volume *coreV1.PersistentVolume, secretName, namespace string) bool {
	// Check if the volume uses a CSI driver
	if volume.Spec.CSI == nil {
		return false
	}

	// Check if the volume uses a NodePublishSecretRef
	if volume.Spec.CSI.NodePublishSecretRef == nil {
		return false
	}

	// Check if the volume uses the specified secret
	if volume.Spec.CSI.NodePublishSecretRef.Name != secretName {
		return false
	}

	// Check if the volume uses the specified namespace
	if volume.Spec.CSI.NodePublishSecretRef.Namespace != namespace {
		return false
	}

	return true
}

func isMaprCSIBased(volume *coreV1.PersistentVolume) bool {
	// Check if the volume is MapR CSI-based
	if volume.Spec.CSI == nil {
		return false
	}

	// Check if the volume is provisioned by one of the MapR CSI provisioners
	for _, provisioner := range maprCSIProvisioners {
		if volume.Spec.CSI.Driver == provisioner {
			return true
		}
	}

	return false
}
