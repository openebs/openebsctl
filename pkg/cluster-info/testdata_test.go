package cluster_info

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// "cspc-operator,cvc-operator,cstor-admission-webhook,openebs-cstor-csi-node,openebs-cstor-csi-controller"
var cspcOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspcOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-1",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cspc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

// "cspc-operator,cvc-operator,cstor-admission-webhook,openebs-cstor-csi-node,openebs-cstor-csi-controller"
var cspcOperatorEvicted = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspcOperatorEvictedPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-8",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cspc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Evicted"},
}

var cvcOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cvcOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-2",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cvc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var cvcOperatorEvicted = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cvcOperatorEvictedPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-9",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cvc-operator", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Evicted"},
}

var cstorAdmissionWebhook = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cstorAdmissionWebhookPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-3",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "cstor-admission-webhook", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsCstorCsiNode = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsCstorCsiNodePOD",
		Namespace:         "openebs",
		UID:               "some-uuid-4",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-cstor-csi-node", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var openebsCstorCsiController = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "openebsCstorCsiControllerPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-5",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-cstor-csi-controller", "openebs.io/version": "2.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndm = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-6",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "ndm", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}

var ndmOperator = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "ndmOperatorPOD",
		Namespace:         "openebs",
		UID:               "some-uuid-7",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{"openebs.io/component-name": "openebs-ndm-operator", "openebs.io/version": "1.1"},
	},
	Status: corev1.PodStatus{Phase: "Running"},
}
