package generate

import (
	"reflect"
	"testing"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	cstorfake "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCSPC(t *testing.T) {
	type args struct {
		c        *client.K8sClient
		nodes    []string
		devs     int
		poolType string
	}
	tests := []struct {
		name    string
		args    args
		want    *cstorv1.CStorPoolCluster
		str     string
		wantErr bool
	}{
		{
			"no cstor installation present",
			args{
				c:     &client.K8sClient{Ns: "", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, &cstorv1.CStorPoolCluster{},
			"", true,
		},
		{
			"stripe kind CSPC with one block-device",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, &cstorv1.CStorPoolCluster{},
			"", false,
		},
		{
			"stripe kind CSPC with two block-device on different nodes",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, &cstorv1.CStorPoolCluster{},
			"", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CSPC(tt.args.c, tt.args.nodes, tt.args.devs, tt.args.poolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSPC() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.str {
				t.Errorf("CSPC() got1 = %v, want %v", got1, tt.str)
			}
		})
	}
}

func Test_isPoolTypeValid(t *testing.T) {
	tests := []struct {
		name      string
		poolNames []string
		want      bool
	}{
		{name: "valid pools", poolNames: []string{"stripe", "mirror", "raidz", "raidz2"}, want: true},
		{name: "invalid pools", poolNames: []string{"striped", "mirrored", "raid-z", "raid-z2", "lvm"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, poolType := range tt.poolNames {
				if got := isPoolTypeValid(poolType); got != tt.want {
					t.Errorf("isPoolTypeValid() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}