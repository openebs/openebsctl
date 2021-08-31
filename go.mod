module github.com/openebs/openebsctl

go 1.16

require (
	github.com/docker/go-units v0.4.0
	github.com/openebs/api/v2 v2.3.0
	github.com/openebs/jiva-operator v1.12.2-0.20210607114402-811a3af7c34a
	github.com/openebs/lvm-localpv v0.6.0
	github.com/openebs/zfs-localpv v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime v0.20.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v1.0.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
