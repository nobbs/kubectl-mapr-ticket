package list

import (
	"fmt"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateTable generates a table from the secrets containing MapR tickets
func GenerateTable(cmd *cobra.Command, items []ListItem) (*metaV1.Table, error) {
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
		Rows: generateRows(items),
	}, nil
}

// generateRows generates the rows for the table from the secrets containing
// MapR tickets
func generateRows(items []ListItem) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(items))

	for _, item := range items {
		row := metaV1.TableRow{
			Object: runtime.RawExtension{
				Object: item.secret,
			},
			Cells: []any{
				item.secret.Name,
				clusterName(item.ticket),
				userName(item.ticket),
				expiryTime(item.ticket),
			},
		}

		rows = append(rows, row)
	}

	return rows
}

// expiryTime returns the expiry time in a human readable format, with an
// indicator if the ticket is expired
func expiryTime(ticket *ticket.MaprTicket) string {
	const timeFormat = time.RFC3339

	if ticket.IsExpired() {
		return fmt.Sprintf("%s (Expired)", ticket.ExpiryTimeToHuman(timeFormat))
	}

	return ticket.ExpiryTimeToHuman(timeFormat)
}

// userName returns the username from the ticket
func userName(ticket *ticket.MaprTicket) string {
	return ticket.UserCreds.GetUserName()
}

// clusterName returns the cluster name from the ticket
func clusterName(ticket *ticket.MaprTicket) string {
	return ticket.Cluster
}
