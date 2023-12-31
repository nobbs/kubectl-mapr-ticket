package list

import (
	"time"

	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	listTableColumns = []metaV1.TableColumnDefinition{
		{
			Name:        "Name",
			Type:        "string",
			Format:      "name",
			Description: "Name of the secret containing the MapR ticket",
			Priority:    0,
		},
		{
			Name:        "MapR Cluster",
			Type:        "string",
			Description: "Name of the MapR cluster that the ticket is for",
			Priority:    0,
		},
		{
			Name:        "User",
			Type:        "string",
			Description: "Name of the MapR user that the ticket is for",
			Priority:    0,
		},
		{
			Name:        "UID",
			Type:        "integer",
			Description: "UID of the MapR user that the ticket is for",
			Priority:    1,
		},
		{
			Name:        "GIDs",
			Type:        "array",
			Description: "GIDs of the MapR user that the ticket is for",
			Priority:    1,
		},
		{
			Name:        "Created",
			Type:        "string",
			Format:      "date-time",
			Description: "Creation time of the ticket",
			Priority:    1,
		},
		{
			Name:        "Expires",
			Type:        "string",
			Format:      "date-time",
			Description: "Expiration time of the ticket",
			Priority:    0,
		},
		{
			Name:        "Status",
			Type:        "string",
			Description: "Status of the ticket",
			Priority:    0,
		},
	}
)

func Print(cmd *cobra.Command, items []ListItem) error {
	switch cmd.Flag("output").Value.String() {
	case "table":
		fallthrough
	case "wide":
		// generate table for output
		table, err := GenerateTable(cmd, items)
		if err != nil {
			return err
		}

		withNamespace := cmd.Flag("all-namespaces").Changed && cmd.Flag("all-namespaces").Value.String() == "true"
		wide := cmd.Flag("output").Value.String() == "wide"

		// print table
		printer := printers.NewTablePrinter(printers.PrintOptions{
			WithNamespace: withNamespace,
			Wide:          wide,
		})

		err = printer.PrintObj(table, cmd.OutOrStdout())
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateTable generates a table from the secrets containing MapR tickets
func GenerateTable(cmd *cobra.Command, items []ListItem) (*metaV1.Table, error) {
	rows := generateRows(items)

	return &metaV1.Table{
		ColumnDefinitions: listTableColumns,
		Rows:              rows,
	}, nil
}

// generateRows generates the rows for the table from the secrets containing
// MapR tickets
func generateRows(items []ListItem) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(items))

	for _, item := range items {
		rows = append(rows, *generateRow(item))
	}

	return rows
}

// generateRow generates a row for the table from the secret containing a MapR
// ticket
func generateRow(item ListItem) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: item.secret,
		},
	}

	var status string
	if item.ticket.IsExpired() {
		status = "Expired"
	} else {
		status = "Valid"
	}

	row.Cells = []any{
		item.secret.Name,
		item.ticket.Cluster,
		item.ticket.UserCreds.GetUserName(),
		item.ticket.UserCreds.GetUid(),
		item.ticket.UserCreds.GetGids(),
		item.ticket.CreateTimeToHuman(time.RFC3339),
		item.ticket.ExpiryTimeToHuman(time.RFC3339),
		status,
	}

	return row
}
