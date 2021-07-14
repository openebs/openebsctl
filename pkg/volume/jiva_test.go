package volume

import (
	openebsFakeClientset "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	jiva "github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"testing"
)

func TestDescribeJivaVolume(t *testing.T) {
	type args struct {
		c   *client.K8sClient
		vol *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeJivaVolume(tt.args.c, tt.args.vol); (err != nil) != tt.wantErr {
				t.Errorf("DescribeJivaVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetJiva(t *testing.T) {

	newScheme := runtime.NewScheme()
	newScheme.AddKnownTypes(schema.GroupVersion{Group: "openebs.io", Version: "v1alpha1"}, &jiva.JivaVolume{})
	scheme.Scheme = newScheme

	type args struct {
		c         *client.K8sClient
		pvList    *corev1.PersistentVolumeList
		openebsNS string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"Test 1",
			args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&jv1, &jv2),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, jivaPV2}},
				openebsNS: "jiva",
			},
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJiva(tt.args.c, tt.args.pvList, tt.args.openebsNS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJiva() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetJiva() got = %v, want %v", got, tt.want)
			}
		})
	}
}
