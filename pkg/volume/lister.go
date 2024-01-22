package volume

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type secretLister interface {
	List() ([]types.TicketSecret, error)
}

type Lister struct {
	client       kubernetes.Interface
	secretLister secretLister

	namespace  string
	secretName string
	sortBy     []SortOptions

	volumes []types.Volume
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
func (l *Lister) List() ([]types.Volume, error) {
	if err := l.getVolumes(); err != nil {
		return nil, err
	}

	l.filterVolumesToMaprCSI().
		filterVolumeUsesTicket().
		collectSecrets().
		sort()

	return l.volumes, nil
}

// getVolumes gets all persistent volumes in the cluster.
func (l *Lister) getVolumes() error {
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	l.volumes = make([]types.Volume, 0, len(pvs.Items))

	for i := range pvs.Items {
		l.volumes = append(
			l.volumes,
			types.Volume{
				Volume: (*types.PersistentVolume)(&pvs.Items[i]),
			},
		)
	}

	return nil
}

// filterVolumesToMaprCSI filters volumes to those that are provisioned by one of the MapR CSI
// provisioners.
func (l *Lister) filterVolumesToMaprCSI() *Lister {
	var filtered []types.Volume

	for _, volume := range l.volumes {
		if volume.Volume.IsMaprCSIBased() {
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
	var filtered []types.Volume

	for _, volume := range l.volumes {
		if volume.Volume.UsesSecret(l.namespace, l.secretName) {
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

	// return early if there are no volumes
	if len(l.volumes) == 0 {
		return l
	}

	// collect all tickets via the secret lister
	secrets, err := l.secretLister.List()
	if err != nil {
		return l
	}

	// add secrets to volumes
	for i := range l.volumes {
		for j := range secrets {
			if l.volumes[i].Volume.Spec.CSI.NodePublishSecretRef.Name == secrets[j].Secret.Name &&
				l.volumes[i].Volume.Spec.CSI.NodePublishSecretRef.Namespace == secrets[j].Secret.Namespace {
				l.volumes[i].Ticket = &secrets[j]
			}
		}
	}

	return l
}
