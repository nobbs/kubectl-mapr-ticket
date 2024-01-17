package secret

import (
	"fmt"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
)

type TicketSecret struct {
	Secret *coreV1.Secret `json:"originalSecret"`
	Ticket *ticket.Ticket `json:"parsedTicket"`
	NumPVC uint32         `json:"numPVC"`
}

func (t *TicketSecret) GetStatus() string {
	if t == nil || t.Ticket == nil {
		return "Invalid"
	}

	if t.Ticket.IsExpired() {
		return fmt.Sprintf("Expired (%s ago)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
	}

	return fmt.Sprintf("Valid (%s left)", util.ShortHumanDurationComparedToNow(t.Ticket.ExpirationTime()))
}
