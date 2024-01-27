// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package secret

import (
	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	tableColumns = []metaV1.TableColumnDefinition{
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
			Name:        "Mapr User",
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

	showInUseTableColumn = metaV1.TableColumnDefinition{
		Name:        "#PVs",
		Type:        "integer",
		Description: "Number of persistent volumes using the ticket",
		Priority:    0,
	}
)

// Print prints the secrets containing MapR tickets in a human-readable format to the given output
// stream.
func Print(cmd *cobra.Command, secrets []types.MaprSecret) error {
	format := cmd.Flag("output").Value.String()
	allNamespaces := cmd.Flag("all-namespaces").Changed && cmd.Flag("all-namespaces").Value.String() == "true"
	withInUse := cmd.Flag("show-in-use").Changed && cmd.Flag("show-in-use").Value.String() == "true"

	// generate table for output
	table := generateTable(secrets)

	// enrich table with in use column
	if withInUse {
		enrichTableWithInUse(table, secrets)
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

	return nil
}

// generateTable generates a table from the secrets containing MapR tickets
func generateTable(secrets []types.MaprSecret) *metaV1.Table {
	rows := generateRows(secrets)

	return &metaV1.Table{
		ColumnDefinitions: tableColumns,
		Rows:              rows,
	}
}

// generateRows generates the rows for the table from the secrets containing
// MapR tickets
func generateRows(secrets []types.MaprSecret) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(secrets))

	for _, item := range secrets {
		rows = append(rows, *generateRow(&item))
	}

	return rows
}

// generateRow generates a row for the table from the secret containing a MapR
// ticket
func generateRow(secrets *types.MaprSecret) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: (*coreV1.Secret)(secrets.Secret),
		},
	}

	row.Cells = []any{
		secrets.Secret.GetName(),
		secrets.Ticket.Cluster,
		secrets.Ticket.UserCreds.GetUserName(),
		secrets.Ticket.UserCreds.GetUid(),
		secrets.Ticket.UserCreds.GetGids(),
		secrets.Ticket.ExpirationTime().Format(ticket.DefaultTimeFormat),
		secrets.GetStatusString(),
		util.ShortHumanDuration(secrets.Ticket.ExpirationTime().Sub(secrets.Ticket.CreationTime())),
		secrets.Ticket.CreationTime().Format(ticket.DefaultTimeFormat),
		util.ShortHumanDurationUntilNow(secrets.Ticket.CreationTime()),
	}

	return row
}

// enrichTableWithInUse enriches the table with a column indicating whether the
// ticket is in use by a persistent volume or not
func enrichTableWithInUse(table *metaV1.Table, secrets []types.MaprSecret) {
	insertPos := len(tableColumns) - 1

	table.ColumnDefinitions = append(
		table.ColumnDefinitions[:insertPos],
		showInUseTableColumn,
		table.ColumnDefinitions[insertPos],
	)

	for i := range table.Rows {
		table.Rows[i].Cells = append(
			table.Rows[i].Cells[:insertPos],
			secrets[i].NumPVC,
			table.Rows[i].Cells[insertPos],
		)
	}
}
