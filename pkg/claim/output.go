// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package claim

import (
	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	tableColumnDefinitions = []metaV1.TableColumnDefinition{
		{
			Name:        "Name",
			Type:        "string",
			Format:      "name",
			Description: "Name of the persistent volume claim",
			Priority:    0,
		},
		{
			Name:        "Secret Namespace",
			Type:        "string",
			Description: "Namespace of the secret containing the MapR ticket",
			Priority:    0,
		},
		{
			Name:        "Secret",
			Type:        "string",
			Description: "Name of the secret containing the MapR ticket",
			Priority:    0,
		},
		{
			Name:        "Volume Name",
			Type:        "string",
			Description: "Name of the persistent volume",
			Priority:    0,
		},
		{
			Name:        "Volume Path",
			Type:        "string",
			Description: "Path of the volume on the MapR cluster",
			Priority:    1,
		},
		{
			Name:        "Volume Handle",
			Type:        "string",
			Description: "Handle of the volume on the MapR cluster",
			Priority:    1,
		},
		{
			Name:        "Ticket Status",
			Type:        "string",
			Description: "Status of the MapR ticket",
			Priority:    0,
		},
		{
			Name:        "Age",
			Type:        "string",
			Format:      "date-time",
			Description: "Creation time of the volume",
			Priority:    0,
		},
	}
)

// Print prints the volume claims to the given output stream in a tabular format known by kubectl.
func Print(cmd *cobra.Command, volumeClaims []types.MaprVolumeClaim) error {
	format := cmd.Flag("output").Value.String()
	allNamespaces := cmd.Flag("all-namespaces").Changed && cmd.Flag("all-namespaces").Value.String() == "true"

	// generate the table
	table := generableTable(volumeClaims)

	// print the table
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

// generableTable generates a table from the given volume claims.
func generableTable(volumeClaims []types.MaprVolumeClaim) *metaV1.Table {
	rows := generateRows(volumeClaims)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

// generateRows generates the rows for the given volume claims.
func generateRows(volumeClaims []types.MaprVolumeClaim) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(volumeClaims))

	for _, pv := range volumeClaims {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

// generateRow generates a row for the given volume claim.
func generateRow(volumeClaim *types.MaprVolumeClaim) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: (*coreV1.PersistentVolumeClaim)(volumeClaim.Claim),
		},
	}

	row.Cells = []any{
		volumeClaim.Claim.GetName(),
		volumeClaim.Volume.GetSecretNamespace(),
		volumeClaim.Volume.GetSecretName(),
		volumeClaim.Volume.GetName(),
		volumeClaim.Volume.GetVolumePath(),
		volumeClaim.Volume.GetVolumeHandle(),
		volumeClaim.Ticket.GetStatusString(),
		util.ShortHumanDurationUntilNow(volumeClaim.Claim.CreationTimestamp.Time),
	}

	return row
}
