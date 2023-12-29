package ticket

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	secretTicketKey = "CONTAINER_TICKET"
)

type List struct {
	client typedV1.SecretInterface
}

func NewList(client typedV1.SecretInterface) *List {
	return &List{
		client: client,
	}
}

func (l *List) Run() ([]coreV1.Secret, error) {
	secrets, err := l.client.List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// filter secrets that don't contain a ticket
	filtered := l.filterSecrets(secrets.Items)

	return filtered, nil
}

func (l *List) filterSecrets(secrets []coreV1.Secret) []coreV1.Secret {
	var filtered []coreV1.Secret

	for _, secret := range secrets {
		if _, ok := secret.Data[secretTicketKey]; ok {
			filtered = append(filtered, secret)
		}
	}

	return filtered
}
