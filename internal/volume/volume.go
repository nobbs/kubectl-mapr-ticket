package volume

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"

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

type secretLister interface {
	List() ([]util.TicketSecret, error)
}

type Lister struct {
	client     kubernetes.Interface
	secretName string
	namespace  string

	secretLister secretLister

	volumes []util.Volume
}

type ListerOption func(*Lister)

func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
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

func (l *Lister) List() ([]util.Volume, error) {
	if err := l.getVolumes(); err != nil {
		return nil, err
	}

	l.filterVolumesToMaprCSI().
		filterVolumeUsesTicket().
		collectSecrets()

	return l.volumes, nil
}

func (l *Lister) getVolumes() error {
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	volumes := make([]util.Volume, 0, len(pvs.Items))

	for i := range pvs.Items {
		volumes = append(volumes, util.Volume{
			PV: &pvs.Items[i],
		})
	}

	l.volumes = volumes

	return nil
}

func (l *Lister) filterVolumesToMaprCSI() *Lister {
	var filtered []util.Volume

	for _, volume := range l.volumes {
		if IsMaprCSIBased(volume.PV) {
			filtered = append(filtered, volume)
		}
	}

	l.volumes = filtered

	return l
}

func (l *Lister) filterVolumeUsesTicket() *Lister {
	var filtered []util.Volume

	for _, volume := range l.volumes {
		if l.TicketUsesSecret(volume.PV, &coreV1.SecretReference{
			Name:      l.secretName,
			Namespace: l.namespace,
		}) {
			filtered = append(filtered, volume)
		}
	}

	l.volumes = filtered

	return l
}

func (l *Lister) TicketUsesSecret(volume *coreV1.PersistentVolume, secretRef *coreV1.SecretReference) bool {
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

func IsMaprCSIBased(volume *coreV1.PersistentVolume) bool {
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

func (l *Lister) collectSecrets() *Lister {
	// check if we have a secret lister, if not, return early
	if l.secretLister == nil {
		return l
	}

	// collect secrets
	secrets, err := l.secretLister.List()
	if err != nil {
		return l
	}

	// add secrets to volumes
	for i := range l.volumes {
		for j := range secrets {
			if l.volumes[i].PV.Spec.CSI.NodePublishSecretRef.Name == secrets[j].Secret.Name {
				l.volumes[i].Ticket = secrets[j].Ticket
			}
		}
	}

	return l
}
