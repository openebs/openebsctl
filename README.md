# openebsctl

This repository is a WIP for the creation of a command line tool plugin for kubectl

# Instructions to build

- clone/download this repo to you system
- run `make openebsctl`
- run `kubectl openebs [get|describe] [resource]` to use the plugin

# Usage


```bash
# Get volumes
$ kubectl openebs get volumes
Namespace  Name                                      Status   Version     Capacity  StorageClass          Attached  Access Mode      Attached Node
---------  ----                                      ------   -------     --------  ------------          --------  -----------      -------------
openebs    pvc-cb978ab8-9045-4d40-abc5-98dfd4fd82fb  Healthy  master-dev  5Gi       cstor.csi.openebs.io  Attached  ReadWriteOnce    vanisingh
openebs    pvc-e20c1212-1ef6-42c4-9638-0145fa3fb4f9  Healthy  master-dev  5Gi       N/A                   N/A                        N/A


# Describe a single volume
$ kubectl openebs describe volume pvc-cbe030cb-63ca-4dfd-ba57-7719a8c93fb2
Volume Details :
----------------
Name            : pvc-cbe030cb-63ca-4dfd-ba57-7719a8c93fb2
Access Mode     : ReadWriteOnce
CSI Driver      : cstor.csi.openebs.io
Storage Class   : openebs-csi-cstor-sparse
Volume Phase    : Bound
Version         : master-dev
CSPC            : cspc-stripe
Size            : 5Gi
Status          : Healthy
ReplicaCount	: 1

Portal Details :
----------------
IQN             :  iqn.2016-09.com.openebs.cstor:pvc-cbe030cb-63ca-4dfd-ba57-7719a8c93fb2
Volume          :  pvc-cbe030cb-63ca-4dfd-ba57-7719a8c93fb2
TargetNodeName  :  vanisingh
Portal          :  10.103.173.0:3260
TargetIP        :  10.103.173.0

Replica Details :
----------------
Name                                                        Pool Instance     Status
----                                                        -------------     ------
pvc-cbe030cb-63ca-4dfd-ba57-7719a8c93fb2-cspc-stripe-56pv   cspc-stripe-56pv  Healthy

# Get CStor pools
$ kubectl openebs get pools
Name              Namespace  HostName                       Free    Capacity   ReadOnly  ProvisionedReplicas  HealthyReplicas  Status  Age
----              ---------  --------                       ----    --------   --------  -------------------  ---------------  ------  ---
fastssd-cstor     test       director-dev-cluster-1-node-1  48200M  48202370k  false     1                    1                ONLINE  2d5h
```
