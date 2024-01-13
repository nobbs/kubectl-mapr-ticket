package volume

import (
	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"

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
			Name:        "Age",
			Type:        "string",
			Format:      "date-time",
			Description: "Creation time of the volume",
			Priority:    0,
		},
	}
)

func Print(cmd *cobra.Command, volumes []coreV1.PersistentVolume) error {
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

func generableTable(pvs []coreV1.PersistentVolume) *metaV1.Table {
	rows := generateRows(pvs)

	return &metaV1.Table{
		ColumnDefinitions: tableColumnDefinitions,
		Rows:              rows,
	}
}

func generateRows(pvs []coreV1.PersistentVolume) []metaV1.TableRow {
	rows := make([]metaV1.TableRow, 0, len(pvs))

	for _, pv := range pvs {
		rows = append(rows, *generateRow(&pv))
	}

	return rows
}

func generateRow(pv *coreV1.PersistentVolume) *metaV1.TableRow {
	row := &metaV1.TableRow{
		Object: runtime.RawExtension{
			Object: pv,
		},
	}

	row.Cells = []any{
		pv.Name,
		getNodePublishSecretRefNamespace(pv),
		getNodePublishSecretRefName(pv),
		getClaimNamespace(pv),
		getClaimName(pv),
		pv.Spec.CSI.VolumeAttributes["volumePath"],
		pv.Spec.CSI.VolumeHandle,
		util.HumanDurationUntilNow(pv.CreationTimestamp.Time),
	}

	return row
}

func getNodePublishSecretRefName(pv *coreV1.PersistentVolume) string {
	if pv.Spec.CSI != nil && pv.Spec.CSI.NodePublishSecretRef != nil {
		return pv.Spec.CSI.NodePublishSecretRef.Name
	}

	return ""
}

func getNodePublishSecretRefNamespace(pv *coreV1.PersistentVolume) string {
	if pv.Spec.CSI != nil && pv.Spec.CSI.NodePublishSecretRef != nil {
		return pv.Spec.CSI.NodePublishSecretRef.Namespace
	}

	return ""
}

func getClaimName(pv *coreV1.PersistentVolume) string {
	if pv.Spec.ClaimRef != nil {
		return pv.Spec.ClaimRef.Name
	}

	return ""
}

func getClaimNamespace(pv *coreV1.PersistentVolume) string {
	if pv.Spec.ClaimRef != nil {
		return pv.Spec.ClaimRef.Namespace
	}

	return ""
}
