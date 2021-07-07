package volume

import (
	"fmt"
	"os"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetLVMLocalPV returns a list of LVM-LocalPV volumes
func GetLVMLocalPV(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	for _, pv := range pvList.Items {
		var attachedNode, storageVersion, customStatus, ns string
		var rows []metav1.TableRow
		_, lvmVolMap, err := c.GetLVMvol(nil, util.Map, "", util.MapOptions{Key: util.Name})
		if err != nil {
			return nil, fmt.Errorf("failed to list ZFSVolumes")
		}
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.LocalPVLVMCSIDriver {
			lvmVol, ok := lvmVolMap[pv.Name]
			if !ok {
				// condition not possible
				_, _ = fmt.Fprintf(os.Stderr, "couldn't find LVM volume "+pv.Name)
			}
			ns = lvmVol.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			accessMode := pv.Spec.AccessModes[0]
			customStatus = lvmVol.Status.State
			storageVersion = ""
			attachedNode = lvmVol.Spec.OwnerNodeID
			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{
					ns, pv.Name, customStatus, storageVersion, pv.Spec.Capacity.Storage(), pv.Spec.StorageClassName, pv.Status.Phase,
					accessMode, attachedNode}})
		}
	}
	return nil, nil
}
