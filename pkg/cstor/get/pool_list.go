/*
Copyright 2020-2021 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package get

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
)

//RunPoolsList fetchs & lists the pools
func RunPoolsList(pools []string, openebsNs string) error {
	k8sClient, err := client.NewK8sClient(openebsNs)
	util.CheckErr(err, util.Fatal)
	if openebsNs == "" {
		nsFromCli, err := k8sClient.GetOpenEBSNamespace(util.CstorCasType)
		if err != nil {
			//return errors.Wrap(err, "Error determining the openebs namespace, please specify using \"--openebs-namespace\" flag")
			return errors.New("no cstor pools found in the cluster")
		}
		k8sClient.Ns = nsFromCli
	}
	var cpools *v1.CStorPoolInstanceList
	if len(pools) == 0 {
		// List all
		cpools, err = k8sClient.GetCSPIs(nil, "")
	} else {
		// Get one or more
		cpools, err = k8sClient.GetCSPIs(pools, "")
	}

	if err != nil {
		return errors.Wrap(err, "error listing pools")
	}
	var rows []metav1.TableRow
	for _, item := range cpools.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{
			item.ObjectMeta.Name,
			item.ObjectMeta.Labels["kubernetes.io/hostname"],
			util.ConvertToIBytes(item.Status.Capacity.Free.String()),
			util.ConvertToIBytes(item.Status.Capacity.Total.String()),
			item.Status.ReadOnly,
			item.Status.ProvisionedReplicas,
			item.Status.HealthyReplicas,
			item.Status.Phase,
			util.Duration(time.Since(item.ObjectMeta.CreationTimestamp.Time))}})
	}
	if len(cpools.Items) == 0 {
		fmt.Println("No Pools are found")
	} else {
		util.TablePrinter(util.CstorPoolListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	}
	return nil
}
