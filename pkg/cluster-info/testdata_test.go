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

package cluster_info

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cspcOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspcOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cspc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var cspcOperatorEvicted = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspcOperatorEvictedPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-8",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cspc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Evicted"},
}

var cvcOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cvcOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-2",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cvc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var cvcOperatorEvicted = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cvcOperatorEvictedPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-9",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cvc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Evicted"},
}

var cstorAdmissionWebhook = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cstorAdmissionWebhookPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-3",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cstor-admission-webhook", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsCstorCsiNode = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsCstorCsiNodePOD",
		Namespace:         "openebs",
		UID:               "some-uuid-4",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-cstor-csi-node", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsCstorCsiController = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsCstorCsiControllerPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-5",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-cstor-csi-controller", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndm = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-6",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "ndm", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndmOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-7",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-ndm-operator", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var localpvProvisionerInOpenebs = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "localpvprovisionerInOpenebsPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-10",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-localpv-provisioner", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var localpvProvisioner = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "localpvprovisionerPOD",
		Namespace:         "xyz",
		UID:               "some-uuid-10",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-localpv-provisioner", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsJivaCsiNode = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsJivaCsiNodePOD",
		Namespace:         "openebs",
		UID:               "some-uuid-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-jiva-csi-node", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var jivaOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "jivaOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-8",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "jiva-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsJivaCsiController = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsJivaCsiControllerPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-2",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-jiva-csi-controller", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndmXYZ = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmPOD",
		Namespace:         "xyz",
		UID:               "some-uuid-6",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "ndm", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndmOperatorXYZ = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmOperatorPOD",
		Namespace:         "xyz",
		UID:               "some-uuid-7",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-ndm-operator", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}
