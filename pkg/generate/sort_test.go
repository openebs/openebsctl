package generate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/api/v2/pkg/apis/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerate(t *testing.T) {
	type args struct {
		list v1alpha1.BlockDeviceList
	}
	tests := []struct {
		name string
		args args
		want *DeviceList
	}{
		{"empty node LinkedList",
			args{list: v1alpha1.BlockDeviceList{Items: []v1alpha1.BlockDevice{}}}, nil},
		{"single node LinkedList",
			args{list: v1alpha1.BlockDeviceList{Items: []v1alpha1.BlockDevice{goodBD1N1}}},
			&DeviceList{goodBD1N1, nil},
		},
		{
			"two node LinkedList",
			args{list: v1alpha1.BlockDeviceList{Items: []v1alpha1.BlockDevice{goodBD1N1, goodBD1N2}}},
			&DeviceList{goodBD1N1, &DeviceList{goodBD1N2, nil}},
		},
		{
			"four node LinkedList",
			args{list: v1alpha1.BlockDeviceList{Items: []v1alpha1.BlockDevice{goodBD1N1, goodBD1N2, goodBD1N3, goodBD2N1}}},
			&DeviceList{goodBD1N1, &DeviceList{goodBD1N2, &DeviceList{goodBD1N3,
				&DeviceList{goodBD2N1, nil}}}},
		},
		{
			"five node LinkedList",
			args{list: v1alpha1.BlockDeviceList{Items: []v1alpha1.BlockDevice{goodBD1N1, goodBD1N2, goodBD1N3, goodBD2N1, goodBD3N1}}},
			&DeviceList{goodBD1N1, &DeviceList{goodBD1N2, &DeviceList{goodBD1N3,
				&DeviceList{goodBD2N1, &DeviceList{goodBD3N1, nil}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Generate(tt.args.list), "Generate(%v)", tt.args.list)
		})
	}
}

func TestDeviceList_Select(t *testing.T) {
	type args struct {
		head  *DeviceList
		size  resource.Quantity
		count int
	}
	tests := []struct {
		name    string
		args    args
		want    []v1alpha1.BlockDevice
		wantErr bool
	}{
		{"empty node LinkedList", args{nil, resource.MustParse("1Gi"), 1}, nil, true},
		{"one node LinkedList", args{&DeviceList{goodBD1N1, nil}, resource.MustParse("0Ki"), 1}, nil, true},
		{"single node LinkedList", args{&DeviceList{goodBD1N1, nil}, resource.MustParse("1Gi"), 1},
			[]v1alpha1.BlockDevice{goodBD1N1}, false},
		{"two node LinkedList, one BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1, nil}},
			resource.MustParse("1Gi"), 1}, []v1alpha1.BlockDevice{goodBD1N1}, false},
		{"two node LinkedList, two BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1, nil}},
			resource.MustParse("1Gi"), 2}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1}, false},
		{"three node LinkedList, one BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1, &DeviceList{goodBD3N1, nil}}},
			resource.MustParse("1Gi"), 1}, []v1alpha1.BlockDevice{goodBD1N1}, false},
		{"three node LinkedList, two BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1, &DeviceList{goodBD3N1, nil}}},
			resource.MustParse("1Gi"), 2}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1}, false},
		{"three node LinkedList, three BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1, &DeviceList{goodBD3N1, nil}}},
			resource.MustParse("1Gi"), 3}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1}, false},
		{"four node LinkedList, four BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, nil}}}},
			resource.MustParse("1Gi"), 4}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1}, false},
		{"four node LinkedList, three BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, nil}}}},
			resource.MustParse("1Gi"), 3}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1}, false},
		{"five node LinkedList, five BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, &DeviceList{goodBD5N1, nil}}}}},
			resource.MustParse("1Gi"), 5}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1, goodBD5N1}, false},
		{"six node LinkedList, four BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, &DeviceList{goodBD5N1, &DeviceList{goodBD6N1, nil}}}}}},
			resource.MustParse("1Gi"), 4}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1}, false},
		{"six node LinkedList, five BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, &DeviceList{goodBD5N1, &DeviceList{goodBD6N1, nil}}}}}},
			resource.MustParse("1Gi"), 5}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1, goodBD5N1}, false},
		{"six node LinkedList, six BD required", args{&DeviceList{goodBD1N1, &DeviceList{goodBD2N1,
			&DeviceList{goodBD3N1, &DeviceList{goodBD4N1, &DeviceList{goodBD5N1, &DeviceList{goodBD6N1, nil}}}}}},
			resource.MustParse("1Gi"), 6}, []v1alpha1.BlockDevice{goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1, goodBD5N1, goodBD6N1}, false},
		{"six node LinkedList, two BD required of 1G", args{bdLinkedList(6, []int{1, 2, 3, 4, 5, 6}), resource.MustParse("1G"), 2}, nil, true},
		{"six node LinkedList, two BD required of 1G", args{bdLinkedList(6, []int{1, 2, 3, 4, 6, 6}), resource.MustParse("1G"), 2},
			[]v1alpha1.BlockDevice{bdGen(5, 6), bdGen(6, 6)}, false},
		{"six node LinkedList, two BD required of 1G", args{bdLinkedList(6, []int{1, 2, 4, 4, 6, 6}), resource.MustParse("1G"), 2},
			[]v1alpha1.BlockDevice{bdGen(3, 4), bdGen(4, 4)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Select(tt.args.head, tt.args.size, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Select() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = got
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select(...), got %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func bdGen(bdSuffix int, GBsize int) v1alpha1.BlockDevice {
	parse := resource.MustParse(fmt.Sprintf("%d", GBsize) + "G")
	return v1alpha1.BlockDevice{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BlockDevice",
			APIVersion: "openebs.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("bd-%d", bdSuffix),
			Namespace: "openebs",
			Labels:    map[string]string{types.HostNameLabelKey: "node-X"},
		},
		Spec: v1alpha1.DeviceSpec{
			Capacity:       v1alpha1.DeviceCapacity{Storage: uint64(parse.Value())},
			NodeAttributes: v1alpha1.NodeAttribute{NodeName: "node-X"},
		},
		Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive},
	}
}

func bdLinkedList(limit int, size []int) *DeviceList {
	if len(size) != limit {
		return nil
	}
	head := &DeviceList{bdGen(1, size[0]), nil}
	curr := head
	for i := 2; i <= limit; i++ {
		curr.next = New(bdGen(i, size[i-1]))
	}
	return head
}
