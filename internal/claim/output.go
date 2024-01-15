package claim

import (
	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"

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
			Name:        "Age",
			Type:        "string",
			Format:      "date-time",
			Description: "Creation time of the volume",
			Priority:    0,
		},
	}
)

func Print(cmd *cobra.Command, volumeClaims []VolumeClaim) error {
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

func generableTable(volumeClaims []VolumeClaim) *metaV1.Table {
	rows := generateRows(volumeClaims)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

func generateRows(volumeClaims []VolumeClaim) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(volumeClaims))

	for _, pv := range volumeClaims {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

func generateRow(volumeClaims *VolumeClaim) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: volumeClaims.pvc,
		},
	}

	row.Cells = []any{
		volumeClaims.pvc.Name,
		getNodePublishSecretRefNamespace(volumeClaims),
		getNodePublishSecretRefName(volumeClaims),
		util.HumanDurationUntilNow(volumeClaims.pvc.CreationTimestamp.Time),
	}

	return row
}

func getNodePublishSecretRefName(volumeClaim *VolumeClaim) string {
	if volumeClaim.pv.Spec.CSI != nil && volumeClaim.pv.Spec.CSI.NodePublishSecretRef != nil {
		return volumeClaim.pv.Spec.CSI.NodePublishSecretRef.Name
	}

	return ""
}

func getNodePublishSecretRefNamespace(volumeClaim *VolumeClaim) string {
	if volumeClaim.pv.Spec.CSI != nil && volumeClaim.pv.Spec.CSI.NodePublishSecretRef != nil {
		return volumeClaim.pv.Spec.CSI.NodePublishSecretRef.Namespace
	}

	return ""
}