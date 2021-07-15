package storage

import (
	"testing"

	fakelvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/fake"
	"github.com/openebs/openebsctl/pkg/client"
)

func TestGetVolumeGroup(t *testing.T) {
	type args struct {
		c  *client.K8sClient
		vg []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no LVM volumegroups present",
			args{
				c: &client.K8sClient{
					Ns:        "lvmlocalpv",
					K8sCS:     nil,
					OpenebsCS: nil,
					LVMCS:     fakelvmclient.NewSimpleClientset(),
				},
				vg: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetVolumeGroups(tt.args.c, tt.args.vg); (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
