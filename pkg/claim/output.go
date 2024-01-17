package claim

import (
	"github.com/spf13/cobra"

	apiClaim "github.com/nobbs/kubectl-mapr-ticket/pkg/api/claim"
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
			Priority:    1,
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

func Print(cmd *cobra.Command, volumeClaims []apiClaim.VolumeClaim) error {
	format := cmd.Flag("output").Value.String()

	// generate the table
	table := generableTable(volumeClaims)

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

func generableTable(volumeClaims []apiClaim.VolumeClaim) *metaV1.Table {
	rows := generateRows(volumeClaims)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

func generateRows(volumeClaims []apiClaim.VolumeClaim) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(volumeClaims))

	for _, pv := range volumeClaims {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

func generateRow(volumeClaim *apiClaim.VolumeClaim) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: volumeClaim.Claim,
		},
	}

	volume := apiVolume.NewVolume(volumeClaim.Volume)

	row.Cells = []any{
		volumeClaim.Claim.Name,
		volume.SecretNamespace(),
		volume.SecretName(),
		volume.Name(),
		volume.VolumePath(),
		volume.VolumeHandle(),
		volumeClaim.Ticket.GetStatus(),
		util.HumanDurationUntilNow(volumeClaim.Claim.CreationTimestamp.Time),
	}

	return row
}
