package volume

import (
	"fmt"
	"os"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Jiva volume methods
type Jiva struct {
	// Volumes are a list of PersistentVolumes which may or may not be CStor CSI provisioned
	Volumes *corev1.PersistentVolumeList
	// k8sClient is the k8sClient to fetch
	k8sClient *client.K8sClient
	// properties like cas-type, filter
	properties map[string]string
}

// Get returns a list of JivaVolumes
func (j *Jiva) Get() ([]metav1.TableRow, error) {
	pvList := j.Volumes
	// 2. Fetch all relevant volume CRs without worrying about openebsNS
	jvMap, _ := j.k8sClient.GetJivaVolumeMap()
	openebsNS := j.properties["openebs-ns"]
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		name := pv.Name
		capacity := pv.Spec.Capacity.Storage()
		sc := pv.Spec.StorageClassName
		attached := pv.Status.Phase
		var attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.JivaCSIDriver {
			jv, ok := jvMap[pv.Name]
			if !ok {
				_, _ = fmt.Fprintln(os.Stderr, "couldn't find jv "+pv.Name)
			}
			ns = jv.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			customStatus = jv.Status.Status // RW, RO, etc
			attachedNode = jv.Labels["nodeID"]
			storageVersion = jv.VersionDetails.Status.Current
		} else {
			// Skip non-CStor & non-Jiva options
			continue
		}
		accessMode := pv.Spec.AccessModes[0]
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				ns, name, customStatus, storageVersion, capacity, sc, attached,
				accessMode, attachedNode},
		})
	}
	return rows, nil
}
