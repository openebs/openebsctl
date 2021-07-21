package persistentvolumeclaim

import (
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestDescribeJivaVolumeClaim(t *testing.T) {
	type args struct {
		c   *client.K8sClient
		pvc *corev1.PersistentVolumeClaim
		vol *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeJivaVolumeClaim(tt.args.c, tt.args.pvc, tt.args.vol); (err != nil) != tt.wantErr {
				t.Errorf("DescribeJivaVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
