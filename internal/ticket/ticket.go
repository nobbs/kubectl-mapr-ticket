package ticket

import (
	"fmt"
	"time"

	"github.com/nobbs/mapr-ticket-parser/pkg/parse"
	coreV1 "k8s.io/api/core/v1"
)

// Wrapper around parse.MaprTicket to add methods
type MaprTicket parse.MaprTicket

// NewTicketFromSecret parses the ticket from the secret and returns it
func NewTicketFromSecret(secret *coreV1.Secret) (*MaprTicket, error) {
	// get ticket from secret
	ticketBytes, ok := secret.Data[secretTicketKey]
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
func (ticket *MaprTicket) isExpired() bool {
	t := time.Unix(int64(ticket.GetExpiryTime()), 0)
	return time.Now().Before(t)
}

// expiryTimeToHuman returns the expiry time in a human readable format
func (ticket *MaprTicket) expiryTimeToHuman(format string) string {
	t := time.Unix(int64(ticket.GetExpiryTime()), 0)
	return t.Format(format)
}
