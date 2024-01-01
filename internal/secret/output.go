package secret

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/yaml"
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
			Name:        "Expiry Time",
			Type:        "string",
			Format:      "date-time",
			Description: "Timestamp of the ticket expiry",
			Priority:    1,
		},
		{
			Name:        "Status",
			Type:        "string",
			Description: "Status of the ticket",
			Priority:    0,
		},
		{
			Name:        "Span",
			Type:        "string",
			Description: "Duration of the ticket",
			Priority:    1,
		},
		{
			Name:        "Creation Time",
			Type:        "string",
			Format:      "date-time",
			Description: "Creation time of the ticket",
			Priority:    1,
		},
		{
			Name:        "Age",
			Type:        "string",
			Description: "Time since the ticket was created",
			Priority:    0,
		},
	}
)

func Print(cmd *cobra.Command, items []ListItem) error {
	format := cmd.Flag("output").Value.String()
	allNamespaces := cmd.Flag("all-namespaces").Changed && cmd.Flag("all-namespaces").Value.String() == "true"
	withInUse := cmd.Flag("show-in-use").Changed && cmd.Flag("show-in-use").Value.String() == "true"

	switch format {
	case "table", "wide":
		// generate table for output
		table := generateTable(items)

		// enrich table with in use column
		if withInUse {
			enrichTableWithInUse(table, items)
		}

		// print table
		printer := printers.NewTablePrinter(printers.PrintOptions{
			WithNamespace: allNamespaces,
			Wide:          format == "wide",
		})

		err := printer.PrintObj(table, cmd.OutOrStdout())
		if err != nil {
			return err
		}
	case "json", "yaml":
		err := printEncoded(items, format, cmd.OutOrStdout())
		if err != nil {
			return err
		}
	}

	return nil
}

// generateTable generates a table from the secrets containing MapR tickets
func generateTable(items []ListItem) *metaV1.Table {
	rows := generateRows(items)

	return &metaV1.Table{
		ColumnDefinitions: listTableColumns,
		Rows:              rows,
	}
}

// generateRows generates the rows for the table from the secrets containing
// MapR tickets
func generateRows(items []ListItem) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(items))

	for _, item := range items {
		rows = append(rows, *generateRow(&item))
	}

	return rows
}

// generateRow generates a row for the table from the secret containing a MapR
// ticket
func generateRow(item *ListItem) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: item.Secret,
		},
	}

	row.Cells = []any{
		item.Secret.Name,
		item.Ticket.Cluster,
		item.Ticket.UserCreds.GetUserName(),
		item.Ticket.UserCreds.GetUid(),
		item.Ticket.UserCreds.GetGids(),
		item.Ticket.ExpiryTimeToHuman(time.RFC3339),
		getStatus(item.Ticket),
		util.ShortHumanDuration(item.Ticket.ExpiryTime().Sub(item.Ticket.CreationTime())),
		item.Ticket.CreateTimeToHuman(time.RFC3339),
		util.ShortHumanDurationUntilNow(item.Ticket.CreationTime()),
	}

	return row
}

// enrichTableWithInUse enriches the table with a column indicating whether the
// ticket is in use by a persistent volume or not
func enrichTableWithInUse(table *metaV1.Table, items []ListItem) {
	numColumns := len(listTableColumns)

	table.ColumnDefinitions = append(
		table.ColumnDefinitions[:numColumns-1],
		metaV1.TableColumnDefinition{
			Name:        "In Use",
			Type:        "boolean",
			Description: "Whether the ticket is in use by a persistent volume or not",
			Priority:    0,
		},
		table.ColumnDefinitions[numColumns-1],
	)

	for i := range table.Rows {
		table.Rows[i].Cells = append(
			table.Rows[i].Cells[:numColumns-1],
			items[i].InUse,
			table.Rows[i].Cells[numColumns-1],
		)
	}
}

func printEncoded(items []ListItem, format string, stream io.Writer) error {
	bytesBuffer := bytes.NewBuffer([]byte{})

	if len(items) == 1 {
		// encode single item
		_, err := bytesBuffer.Write(encodeItem(&items[0], format))
		if err != nil {
			return err
		}
	} else {
		// encode multiple items
		_, err := bytesBuffer.Write(encodeItems(items, format))
		if err != nil {
			return err
		}
	}

	// print encoded items
	_, err := bytesBuffer.WriteTo(stream)
	if err != nil {
		return err
	}

	return nil
}

func encodeItems(items []ListItem, format string) []byte {
	switch format {
	case "json":
		encoded, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return nil
		}

		return encoded
	case "yaml":
		encoded, err := yaml.Marshal(items)
		if err != nil {
			return nil
		}

		return encoded
	}

	return nil
}

func encodeItem(item *ListItem, format string) []byte {
	switch format {
	case "json":
		encoded, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return nil
		}

		return encoded
	case "yaml":
		encoded, err := yaml.Marshal(item)
		if err != nil {
			return nil
		}

		return encoded
	}

	return nil
}

func getStatus(ticket *ticket.MaprTicket) string {
	if ticket.IsExpired() {
		return fmt.Sprintf("Expired (%s ago)", util.ShortHumanDurationComparedToNow(ticket.ExpiryTime()))
	}

	return fmt.Sprintf("Valid (%s left)", util.ShortHumanDurationComparedToNow(ticket.ExpiryTime()))
}
