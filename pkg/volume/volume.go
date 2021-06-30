package volume

import (
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Impl -> Headers
// Get manages various implementations of Volume listing
func Get(vols []string, casType, openebsNS string) error {
	k, _ := client.NewK8sClient("")
	var pvList *corev1.PersistentVolumeList
	if vols == nil {
		pvList, _ = k.GetPVs(nil, "")
	} else {
		pvList, _ = k.GetPVs(vols, "")
	}
	impl := []func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error){GetJiva, GetCStor}

	var rows []metav1.TableRow
	// TODO: Decide if running each 7 cas-type implementations in
	// go-routine will be wise & will still maintain the ordering
	for _, t := range impl {
		jr, err := t(k, pvList, openebsNS)
		if err != nil {
			rows = append(rows, jr...)
		}
	}
	util.TablePrinter(util.VolumeListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
