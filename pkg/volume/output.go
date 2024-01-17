package volume

import (
	"github.com/spf13/cobra"

	apiVolume "github.com/nobbs/kubectl-mapr-ticket/pkg/api/volume"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

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

func Print(cmd *cobra.Command, volumes []apiVolume.Volume) error {
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

func generableTable(volumes []apiVolume.Volume) *metaV1.Table {
	rows := generateRows(volumes)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

func generateRows(volumes []apiVolume.Volume) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(volumes))

	for _, pv := range volumes {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

func generateRow(volume *apiVolume.Volume) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: volume.Volume,
		},
	}

	row.Cells = []any{
		volume.Volume.Name,
		volume.SecretNamespace(),
		volume.SecretName(),
		volume.ClaimNamespace(),
		volume.ClaimName(),
		volume.VolumePath(),
		volume.VolumeHandle(),
		volume.Ticket.GetStatus(),
		util.HumanDurationUntilNow(volume.Volume.CreationTimestamp.Time),
	}

	return row
}
