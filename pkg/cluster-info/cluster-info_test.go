package cluster_info

import "testing"

func TestShowClusterInfo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"Test 1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowClusterInfo(); (err != nil) != tt.wantErr {
				t.Errorf("ShowClusterInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
