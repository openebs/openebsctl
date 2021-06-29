package volume

import (
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

type Volume interface {
	// PV -> CV,CVC,CVA/JV/
	// Implicit arguments=PVList
	Get() ([]metav1.TableRow, error)
	// THINKING_POINT: Do I need to implement a single cas-by-cas filter or a
	// SingleTon filter function which can give me all or one as a map
	// Filter(volumes []string, properties map[string]string) (*corev1.PersistentVolumeList, error)
	// Describe(volumes []string, properties map[string]string) ()
}

// Impl -> Headers

func GetVolumes(vols []string, casType, openebsNS string) error {
	k, _ := client.NewK8sClient("")
	var pvList *corev1.PersistentVolumeList
	if vols == nil {
		pvList, _ = k.GetPVs(nil, "")
	} else {
		pvList, _ = k.GetPVs(vols, "")
	}
	prop := map[string]string{
		"casType":    casType,
		"openebs-ns": openebsNS,
	}
	types := []Volume{&Jiva{
		k8sClient:  k,
		Volumes:    pvList,
		properties: prop,
	}, &CStor{
		k8sClient:  k,
		Volumes:    pvList,
		properties: prop,
	}}
	var rows []metav1.TableRow
	for _, t := range types {
		jr, _ := t.Get()
		rows = append(rows, jr...)
	}
	util.TablePrinter(util.VolumeListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
