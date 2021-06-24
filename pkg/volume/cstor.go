package volume

import (
	"fmt"
	"os"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CStor volume
type CStor struct {
	// Volumes are a list of PersistentVolumes which may or may not be CStor CSI provisioned
	Volumes *corev1.PersistentVolumeList
	// k8sClient is the k8sClient to fetch
	k8sClient *client.K8sClient
	// OpenEBS namespace
	properties map[string]string
}

// Get implements the Volume interface to Get CStor volumes
func (c *CStor) Get() ([]metav1.TableRow, error) {
	pvList := c.Volumes
	openebsNS := c.properties["openebs-ns"]
	var (
		cvMap  map[string]v1.CStorVolume
		cvaMap map[string]v1.CStorVolumeAttachment
	)
	// TODO: What to do if these throw some errors?, errors need to be of two kinds: warnings/errors
	cvMap, _ = c.k8sClient.GetCStorVolumeMap()
	cvaMap, _ = c.k8sClient.GetCStorVolumeAttachmentMap()
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		var attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		// 2. For eligible PVs fetch the custom-resource to add more info
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.CStorCSIDriver {
			// For all CSI CStor PV, there exist a CV
			cv, ok := cvMap[pv.Name]
			if !ok {
				// condition not possible
				_, _ = fmt.Fprintf(os.Stderr, "couldn't find cv "+pv.Name)
			}
			ns = cv.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			customStatus = string(cv.Status.Phase)
			storageVersion = cv.VersionDetails.Status.Current
			cva, cvaOk := cvaMap[pv.Name]
			if cvaOk {
				attachedNode = cva.Labels["nodeID"]
			}
		}
		// TODO: What should be done for multiple AccessModes
		accessMode := pv.Spec.AccessModes[0]
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				ns, pv.Name, customStatus, storageVersion, pv.Spec.Capacity.Storage(), pv.Spec.StorageClassName, pv.Status.Phase,
				accessMode, attachedNode}})
	}
	return rows, nil
}
