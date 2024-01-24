package types

import (
	"fmt"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
)

type Secret coreV1.Secret

type TicketSecret struct {
	Secret *Secret        `json:"originalSecret"`
	Ticket *ticket.Ticket `json:"parsedTicket"`
	NumPVC uint32         `json:"numPVC"`
}

// GetName returns the name of the secret
func (t *TicketSecret) GetName() string {
	if t == nil || t.Secret == nil {
		return ""
	}

	return t.Secret.GetName()
}

// GetNamespace returns the namespace of the secret
func (t *TicketSecret) GetNamespace() string {
	if t == nil || t.Secret == nil {
		return ""
	}

	return t.Secret.GetNamespace()
}

// GetCluster returns the cluster of the ticket
func (t *TicketSecret) GetCluster() string {
	if t == nil || t.Ticket == nil {
		return ""
	}

	return t.Ticket.GetCluster()
}

// GetUser returns the user of the ticket
func (t *TicketSecret) GetUser() string {
	if t == nil || t.Ticket == nil {
		return ""
	}

	return t.Ticket.GetUser()
}

// GetExpirationTime returns the expiration time of the ticket
func (t *TicketSecret) GetExpirationTime() time.Time {
	if t == nil || t.Ticket == nil {
		return time.Time{}
	}

	return t.Ticket.ExpirationTime()
}

// GetCreationTime returns the creation time of the ticket
func (t *TicketSecret) GetCreationTime() time.Time {
	if t == nil || t.Ticket == nil {
		return time.Time{}
	}

	return t.Ticket.CreationTime()
}

// GetStatusString returns a human readable string describing the status of the ticket
func (t *TicketSecret) GetStatusString() string {
	if t == nil {
		return "Not found / Invalid"
	}

	if t.Secret == nil {
		return "No secret found"
	}

	if t == nil || t.Ticket == nil {
		return "No ticket found"
	}

	if t.Ticket.IsExpired() {
		return fmt.Sprintf("Expired (%s ago)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
	}

	return fmt.Sprintf("Valid (%s left)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
}
