<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# LOCALPV-ZFS Storage Engine Commands

## Table of Contents
* [LocalPV-ZFS](#localpv-zfs)
    * [Get LocalPV-ZFS volumes](#get-localpv-zfs-volumes)
    * [Get LocalPV-ZFS Pools](#get-localpv-zfs-pools)
    * [Describe LocalPV-ZFS volumes](#describe-localpv-zfs-volumes)
    * [Describe LocalPV-ZFS pools](#describe-localpv-zfs-pools)
    * [Describe LocalPV-ZFS PVCs](#describe-localpv-zfs-pvcs)
  
* #### `LocalPV-ZFS`
    * #### Get `LocalPV-ZFS` volumes
      ```bash
      $ kubectl openebs get volumes --cas-type=localpv-zfs
      NAMESPACE   NAME                                       STATUS   VERSION   CAPACITY   STORAGE CLASS   ATTACHED   ACCESS MODE     ATTACHED NODE
      openebs     pvc-43fcbc72-a45a-49d5-9ec3-e383fcb91452   Ready    1.9.0     4Gi        openebs-zfspv   Bound      ReadWriteOnce   worker-sh1
      ```
      Note: For volumes not attached to any application, the `ATTACH NODE` would be shown as `N/A`.
    * #### Get `LocalPV-ZFS` Pools
      ```bash
      $ kubectl openebs get storage --cas-type=localpv-zfs
      NAME              FREESIZE
      node1         
      └─zfs-test-pool   32 GiB
      
      node2         
      └─zfs-test-pool   36 GiB
      ```
    * #### Describe `LocalPV-ZFS volumes`
      ```bash
      $ kubectl openebs describe vol pvc-43fcbc72-a45a-49d5-9ec3-e383fcb91452
  
      pvc-43fcbc72-a45a-49d5-9ec3-e383fcb91452 Details :
      -----------------
      Name          : pvc-43fcbc72-a45a-49d5-9ec3-e383fcb91452
      Namespace     : openebs
      AccessMode    : ReadWriteOnce
      CSIDriver     : zfs.csi.openebs.io
      Capacity      : 4Gi
      PVC           : csi-zfspv
      VolumePhase   : Bound
      StorageClass  : openebs-zfspv
      Version       : N/A
      Status        : Ready
      VolumeType    : DATASET
      PoolName      : zfspv-pool
      FileSystem    : zfs
      Compression   : off
      Deduplication : off
      NodeID        : worker-sh1
      Recordsize    : 4k
      ```
    * #### Describe `LocalPV-ZFS` pools
      ```bash
      $ kubectl openebs describe storage node2
    
       node2 Details :
    
       HOSTNAME        : node2
       NAMESPACE       : openebs
       NUMBER OF POOLS : 1
       TOTAL FREE      : 32 GiB
      ```
    * #### Describe `LocalPV-ZFS` PVCs
      ```bash
      $ kubectl openebs describe pvc csi-zfspv

      csi-zfspv Details  :
      -------------------
      NAME               : csi-zfspv
      NAMESPACE          : default
      CAS TYPE           : localpv-zfs
      BOUND VOLUME       : pvc-4aee0a2a-dccd-456b-95af-693ac8108be1
      STORAGE CLASS      : openebs-zfspv
      SIZE               : 4Gi
      PVC STATUS         : Bound

      pvc-4aee0a2a-dccd-456b-95af-693ac8108be1 Details :
      ------------------
      NAME              : pvc-4aee0a2a-dccd-456b-95af-693ac8108be1
      NAMESPACE         : localpv-zfs
      ACCESS MODE       : ReadWriteOnce
      CSI DRIVER        : zfs.csi.openebs.io
      CAPACITY          : 4Gi
      PVC NAME          : csi-zfspv
      VOLUME PHASE      : Bound
      STORAGE CLASS     : openebs-zfspv
      VERSION           : 1.9.1
      ZFS VOLUME STATUS : Ready
      VOLUME TYPE       : DATASET
      POOL NAME         : zfspv-pool
      FILE SYSTEM       : zfs
      COMPRESSION       : off
      DEDUPLICATION     : off
      NODE ID           : node2-virtual-machine
      RECORD SIZE       : 4k
      ```