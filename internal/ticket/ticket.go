package ticket

import (
	"fmt"
	"time"

	"github.com/nobbs/mapr-ticket-parser/pkg/parse"
	coreV1 "k8s.io/api/core/v1"
)

const (
	// SecretMaprTicketKey is the key used for MapR tickets in secrets
	SecretMaprTicketKey = "CONTAINER_TICKET"

	// DefaultTimeFormat is the default time format used for human readable time
	// strings
	DefaultTimeFormat = time.RFC3339
)

// SecretContainsMaprTicket returns true if the secret contains the key typically
// used for MapR tickets
type ErrSecretDoesNotContainMaprTicket struct {
	Name      string
	Namespace string
}

// NewErrSecretDoesNotContainMaprTicket returns a new ErrSecretDoesNotContainMaprTicket
func NewErrSecretDoesNotContainMaprTicket(namespace, name string) ErrSecretDoesNotContainMaprTicket {
	return ErrSecretDoesNotContainMaprTicket{
		Name:      name,
		Namespace: namespace,
	}
}

// Error returns the error message for ErrSecretDoesNotContainMaprTicket
func (err ErrSecretDoesNotContainMaprTicket) Error() string {
	return fmt.Sprintf("secret %s/%s does not contain a MapR ticket", err.Namespace, err.Name)
}

// SecretContainsMaprTicket returns true if the secret contains the key typically
// used for MapR tickets
func SecretContainsMaprTicket(secret *coreV1.Secret) bool {
	_, ok := secret.Data[SecretMaprTicketKey]
	return ok
}

// Wrapper around parse.MaprTicket to add methods
type MaprTicket parse.MaprTicket

// NewMaprTicket returns a new empty MaprTicket
func NewMaprTicket() *MaprTicket {
	return (*MaprTicket)(parse.NewMaprTicket())
}

// NewMaprTicketFromSecret parses the ticket from the secret and returns it
func NewMaprTicketFromSecret(secret *coreV1.Secret) (*MaprTicket, error) {
	// get ticket from secret
	ticketBytes, ok := secret.Data[SecretMaprTicketKey]
	if !ok {
		return nil, NewErrSecretDoesNotContainMaprTicket(secret.Namespace, secret.Name)
	}

	// parse ticket
	ticket, err := parse.Unmarshal(ticketBytes)
	if err != nil {
		return nil, err
	}

	return (*MaprTicket)(ticket), nil
}

// isExpired returns true if the ticket is expired
func (ticket *MaprTicket) IsExpired() bool {
	return time.Now().After(ticket.ExpirationTime())
}

// ExpirationTime returns the expiry time of the ticket as a time.Time object
func (ticket *MaprTicket) ExpirationTime() time.Time {
	return time.Unix(int64(ticket.GetExpiryTime()), 0)
}

// CreationTime returns the creation time of the ticket as a time.Time object
func (ticket *MaprTicket) CreationTime() time.Time {
	return time.Unix(int64(ticket.GetCreationTimeSec()), 0)
}
