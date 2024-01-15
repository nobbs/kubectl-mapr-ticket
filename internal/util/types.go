package util

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"

	coreV1 "k8s.io/api/core/v1"
)

type TicketSecret struct {
	Secret *coreV1.Secret `json:"originalSecret"`
	Ticket *ticket.Ticket `json:"parsedTicket"`
	NumPVC uint32         `json:"numPVC"`
}

type Volume struct {
	PV     *coreV1.PersistentVolume
	Ticket *ticket.Ticket
}

type VolumeClaim struct {
	PVC *coreV1.PersistentVolumeClaim
	PV  *coreV1.PersistentVolume
}
