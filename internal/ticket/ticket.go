package ticket

import (
	"fmt"
	"time"

	"github.com/nobbs/mapr-ticket-parser/pkg/parse"
	coreV1 "k8s.io/api/core/v1"
)

const (
	secretMaprTicketKey = "CONTAINER_TICKET"
)

// SecretContainsMaprTicket returns true if the secret contains the key typically
// used for MapR tickets
func SecretContainsMaprTicket(secret *coreV1.Secret) bool {
	_, ok := secret.Data[secretMaprTicketKey]
	return ok
}

// Wrapper around parse.MaprTicket to add methods
type MaprTicket parse.MaprTicket

// NewTicketFromSecret parses the ticket from the secret and returns it
func NewTicketFromSecret(secret *coreV1.Secret) (*MaprTicket, error) {
	// get ticket from secret
	ticketBytes, ok := secret.Data[secretMaprTicketKey]
	if !ok {
		return nil, fmt.Errorf("secret %s does not contain a MapR ticket", secret.Name)
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
	return time.Now().After(ticket.ExpiryTime())
}

// expiryTimeToHuman returns the expiry time in a human readable format
func (ticket *MaprTicket) ExpiryTimeToHuman(format string) string {
	return ticket.ExpiryTime().Format(format)
}

// createTimeToHuman returns the creation time in a human readable format
func (ticket *MaprTicket) CreateTimeToHuman(format string) string {
	return ticket.CreationTime().Format(format)
}

// ExpiryTime returns the expiry time of the ticket as a time.Time object
func (ticket *MaprTicket) ExpiryTime() time.Time {
	return time.Unix(int64(ticket.GetExpiryTime()), 0)
}

// CreationTime returns the creation time of the ticket as a time.Time object
func (ticket *MaprTicket) CreationTime() time.Time {
	return time.Unix(int64(ticket.GetCreationTimeSec()), 0)
}
