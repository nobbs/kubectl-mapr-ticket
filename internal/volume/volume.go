package volume

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// SecretAll is a special value that can be used to specify that all secrets from the specified
	// namespace should be used.
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
	client       kubernetes.Interface
	secretLister secretLister
	secretName   string
	namespace    string
	sortBy       []SortOptions
	volumes      []util.Volume
}

type ListerOption func(*Lister)

// WithSortBy sets the sort order used by the Lister for output
func WithSortBy(sortBy []SortOptions) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}

// WithSecretLister sets the secret lister used by the Lister to collect secrets and tickets
// referenced by the volumes
func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
	}
}

// NewLister returns a new volume lister that lists volumes that are provisioned by one of the
// MapR CSI provisioners and that use the specified secret.
func NewLister(client kubernetes.Interface, secretName string, namespace string, opts ...ListerOption) *Lister {
	l := &Lister{
		client:     client,
		secretName: secretName,
		namespace:  namespace,
		sortBy:     DefaultSortBy,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// List returns a list of volumes using the MapR CSI provisioners and the specified secret.
func (l *Lister) List() ([]util.Volume, error) {
	if err := l.getVolumes(); err != nil {
		return nil, err
	}

	l.filterVolumesToMaprCSI().
		filterVolumeUsesTicket().
		collectSecrets().
		sort()

	return l.volumes, nil
}

// TicketUsesSecret returns true if the volume uses the specified secret and false otherwise. If the
// secret name is equal to the value of SecretAll, all volumes that use a secret from the specified
// namespace are returned. If the secret namespace is equal to the value of NamespaceAll, basically
// all volumes that any secret in any namespace will evaluate to true.
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
	if secretRef.Namespace == util.NamespaceAll {
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

// getVolumes gets all persistent volumes in the cluster.
func (l *Lister) getVolumes() error {
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	l.volumes = make([]util.Volume, 0, len(pvs.Items))

	for i := range pvs.Items {
		l.volumes = append(l.volumes, util.Volume{
			PV: &pvs.Items[i],
		})
	}

	return nil
}

// filterVolumesToMaprCSI filters volumes to those that are provisioned by one of the MapR CSI
// provisioners.
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

// filterVolumeUsesTicket filters volumes that use the specified ticket secret. If the secret name
// is equal to the value of SecretAll, all volumes that use a secret from the specified namespace
// are returned.
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

// collectSecrets collects secrets and tickets referenced by the volumes, if a secret lister was
// provided to the Lister.
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
			if l.volumes[i].PV.Spec.CSI.NodePublishSecretRef.Name == secrets[j].Secret.Name &&
				l.volumes[i].PV.Spec.CSI.NodePublishSecretRef.Namespace == secrets[j].Secret.Namespace {
				l.volumes[i].Ticket = &secrets[j]
			}
		}
	}

	return l
}

// IsMaprCSIBased returns true if the volume is provisioned by one of the MapR CSI provisioners and
// false otherwise.
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
