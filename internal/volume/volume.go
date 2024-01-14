package volume

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	SecretAll = "<all>"
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

	volumes []coreV1.PersistentVolume
}

func NewLister(client kubernetes.Interface, secretName string, namespace string) *Lister {
	l := &Lister{
		client:     client,
		secretName: secretName,
		namespace:  namespace,
	}

	return l
}

func (l *Lister) List() ([]coreV1.PersistentVolume, error) {
	if err := l.getVolumes(); err != nil {
		return nil, err
	}

	l.filterVolumesToMaprCSI().
		filterVolumeUsesTicket()

	return l.volumes, nil
}

func (l *Lister) getVolumes() error {
	volumes, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	l.volumes = volumes.Items

	return nil
}

func (l *Lister) filterVolumesToMaprCSI() *Lister {
	var filtered []coreV1.PersistentVolume

	for _, volume := range l.volumes {
		if isMaprCSIBased(&volume) {
			filtered = append(filtered, volume)
		}
	}

	l.volumes = filtered

	return l
}

func (l *Lister) filterVolumeUsesTicket() *Lister {
	var filtered []coreV1.PersistentVolume

	for _, volume := range l.volumes {
		if TicketUsesSecret(&volume, &coreV1.SecretReference{
			Name:      l.secretName,
			Namespace: l.namespace,
		}) {
			filtered = append(filtered, volume)
		}
	}

	l.volumes = filtered

	return l
}

func TicketUsesSecret(volume *coreV1.PersistentVolume, secretRef *coreV1.SecretReference) bool {
	// Check if the volume uses a CSI driver
	if volume.Spec.CSI == nil {
		return false
	}

	// Check if the volume uses a NodePublishSecretRef
	if volume.Spec.CSI.NodePublishSecretRef == nil {
		return false
	}

	// Check if we want secrets from all namespaces
	if secretRef.Namespace == metaV1.NamespaceAll {
		return true
	}

	// Check if we want all secrets from the specified namespace
	if secretRef.Name == SecretAll && volume.Spec.CSI.NodePublishSecretRef.Namespace == secretRef.Namespace {
		return true
	}

	// Check if the volume uses the specified secret
	if volume.Spec.CSI.NodePublishSecretRef.Name != secretRef.Name {
		return false
	}

	// Check if the volume uses the specified namespace
	if volume.Spec.CSI.NodePublishSecretRef.Namespace != secretRef.Namespace {
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
