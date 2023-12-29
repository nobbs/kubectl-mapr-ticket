package ticket

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateTable generates a table from the secrets containing MapR tickets
func GenerateTable(cmd *cobra.Command, secrets []coreV1.Secret) (*metaV1.Table, error) {
	return &metaV1.Table{
		ColumnDefinitions: []metaV1.TableColumnDefinition{
			{
				Name: "Name",
			},
			{
				Name: "MapR Cluster",
			},
			{
				Name: "User",
			},
			{
				Name: "Expiration",
			},
		},
		Rows: generateRows(secrets),
	}, nil
}

// generateRows generates the rows for the table from the secrets containing
// MapR tickets
func generateRows(secrets []coreV1.Secret) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(secrets))

	for i := range secrets {
		secret := secrets[i]
		t, err := NewTicketFromSecret(&secret)
		if err != nil {
			continue
		}

		row := metaV1.TableRow{
			Object: runtime.RawExtension{
				Object: &secret,
			},
			Cells: []any{
				secret.Name,
				clusterName(t),
				userName(t),
				expiryTime(t),
			},
		}

		rows = append(rows, row)
	}

	return rows
}

// expiryTime returns the expiry time in a human readable format, with an
// indicator if the ticket is expired
func expiryTime(ticket *MaprTicket) string {
	const timeFormat = time.RFC3339

	if ticket.isExpired() {
		return fmt.Sprintf("%s (Expired)", ticket.expiryTimeToHuman(timeFormat))
	}

	return ticket.expiryTimeToHuman(timeFormat)
}

// userName returns the username from the ticket
func userName(ticket *MaprTicket) string {
	return ticket.UserCreds.GetUserName()
}

// clusterName returns the cluster name from the ticket
func clusterName(ticket *MaprTicket) string {
	return ticket.Cluster
}
