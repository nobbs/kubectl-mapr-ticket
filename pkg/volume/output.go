// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package volume

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
			Description: "Name of the persistent volume",
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
			Name:        "Claim Namespace",
			Type:        "string",
			Description: "Namespace of the persistent volume claim",
			Priority:    0,
		},
		{
			Name:        "Claim",
			Type:        "string",
			Description: "Name of the persistent volume claim",
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

// Print prints the volumes to the default output stream in a human readable table format.
func Print(cmd *cobra.Command, volumes []types.MaprVolume) error {
	format := cmd.Flag("output").Value.String()

	// generate the table
	table := generableTable(volumes)

	// print the table
	printer := printers.NewTablePrinter(printers.PrintOptions{
		Wide: format == "wide",
	})

	err := printer.PrintObj(table, cmd.OutOrStdout())
	if err != nil {
		return err
	}

	return nil
}

// generableTable generates a table from the specified volumes.
func generableTable(volumes []types.MaprVolume) *metaV1.Table {
	rows := generateRows(volumes)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

// generateRows generates the rows for the table from the specified volumes.
func generateRows(volumes []types.MaprVolume) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(volumes))

	for _, pv := range volumes {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

// generateRow generates a row for the table from the specified volume.
func generateRow(volume *types.MaprVolume) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: (*coreV1.PersistentVolume)(volume.Volume),
		},
	}

	row.Cells = []any{
		volume.Volume.GetName(),
		volume.Volume.GetSecretNamespace(),
		volume.Volume.GetSecretName(),
		volume.Volume.GetClaimNamespace(),
		volume.Volume.GetClaimName(),
		volume.Volume.GetVolumePath(),
		volume.Volume.GetVolumeHandle(),
		volume.Ticket.GetStatusString(),
		util.ShortHumanDurationUntilNow(volume.Volume.CreationTimestamp.Time),
	}

	return row
}
