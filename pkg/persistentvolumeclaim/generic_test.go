package persistentvolumeclaim

import (
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestDescribeGenericVolumeClaim(t *testing.T) {
	type args struct {
		pvc     *corev1.PersistentVolumeClaim
		pv      *corev1.PersistentVolume
		casType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All Valid Values",
			args: args{
				pv: &cstorPV1,
				pvc: &cstorPVC1,
				casType: "some-cas",
			},
			wantErr: false,
		},
		{
			name: "PV missing",
			args: args{
				pv: nil,
				pvc: &cstorPVC1,
				casType: "some-cas",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeGenericVolumeClaim(tt.args.pvc, tt.args.pv, tt.args.casType); (err != nil) != tt.wantErr {
				t.Errorf("DescribeGenericVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
