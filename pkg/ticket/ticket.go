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

// Wrapper around parse.Ticket to add methods
type Ticket parse.MaprTicket

// NewMaprTicket returns a new empty MaprTicket
func NewMaprTicket() *Ticket {
	return (*Ticket)(parse.NewMaprTicket())
}

// NewMaprTicketFromSecret parses the ticket from the secret and returns it
func NewMaprTicketFromSecret(secret *coreV1.Secret) (*Ticket, error) {
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

	return (*Ticket)(ticket), nil
}

// GetCluster returns the cluster that the ticket is for
func (ticket *Ticket) GetCluster() string {
	if ticket == nil {
		return ""
	}

	return ticket.Cluster
}

// GetUser returns the user that the ticket is for
func (ticket *Ticket) GetUser() string {
	if ticket == nil || ticket.UserCreds == nil || ticket.UserCreds.UserName == nil {
		return ""
	}

	return *ticket.UserCreds.UserName
}

// isExpired returns true if the ticket is expired
func (ticket *Ticket) IsExpired() bool {
	return time.Now().After(ticket.ExpirationTime())
}

// ExpirationTime returns the expiry time of the ticket as a time.Time object
func (ticket *Ticket) ExpirationTime() time.Time {
	return time.Unix(int64(ticket.GetExpiryTime()), 0)
}

// CreationTime returns the creation time of the ticket as a time.Time object
func (ticket *Ticket) CreationTime() time.Time {
	return time.Unix(int64(ticket.GetCreationTimeSec()), 0)
}

// ExpiresBefore returns true if the ticket expires before the given duration
func (ticket *Ticket) ExpiresBefore(duration time.Duration) bool {
	return ticket.ExpirationTime().Before(time.Now().Add(duration))
}

// AsMaprTicket returns the ticket as a parse.MaprTicket object
func (ticket *Ticket) AsMaprTicket() *parse.MaprTicket {
	return (*parse.MaprTicket)(ticket)
}
