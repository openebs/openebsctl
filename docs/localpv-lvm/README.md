<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# LOCALPV-LVM Storage Engine Commands

## Table of Contents
* [LocalPV-LVM](#localpv-lvm)
    * [Get LocalPV-LVM volumes](#get-localpv-lvm-volumes)
    * [Get LocalPV-LVM VolumeGroups](#get-localpv-lvm-volumegroups)
    * [Describe LocalPV-LVM volumeGroups](#describe-localpv-lvm-volumeGroups)
    * [Describe LocalPV-LVM volumes](#describe-localpv-lvm-volumes)

* #### `LocalPV-LVM`
    * #### Get `LocalPV-LVM` volumes
      ```bash
      $ kubectl openebs get volumes --cas-type=localpv-lvm
      NAMESPACE   NAME                                       STATUS   VERSION   CAPACITY   STORAGE CLASS   ATTACHED   ACCESS MODE     ATTACHED NODE
      openebs     pvc-04c2d4ea-f072-4e17-9e0a-db0fde0b2550   Ready    ci        1Gi        lvmpv-sc        Bound      ReadWriteOnce   worker-sh1
      openebs     pvc-1ec1c9b7-b74e-4742-901d-2af4558d6636   Ready    ci        1Gi        openebs-lvmpv   Bound      ReadWriteOnce   worker-sh1
      openebs     pvc-9999274f-ad01-48bc-9b21-7c51b47a870c   Ready    ci        4Gi        openebs-lvmpv   Bound      ReadWriteOnce   worker-sh1
      ```
      Note: For volumes not attached to any application, the `ATTACH NODE` would be shown as `N/A`.
    * #### Get `LocalPV-LVM` VolumeGroups
      ```bash
      $ kubectl openebs get storage --cas-type=lvmlocalpv
      NAME         FREESIZE   TOTALSIZE
      worker-sh1              
      └─lvmvg      1020 GiB   1024 GiB
      
      worker-sh2              
      └─lvmvg-1    46.7 GiB   50 GiB
      ```
    * #### Describe `LocalPV-LVM` volumeGroups
      ```bash
      $ kubectl openebs describe storage worker-sh1
      worker-sh1 Details :
    
      HOSTNAME        : worker-sh1
      NAMESPACE       : openebs
      NUMBER OF POOLS : 1
      TOTAL CAPACITY  : 1024.0GiB
      TOTAL FREE      : 1020.0GiB
      TOTAL LVs       : 1
      TOTAL PVs       : 1
    
      Volume group details
      ---------------------
      NAME    UUID                                     LV COUNT   PV COUNT   USED PERCENTAGE
      lvmvg   IgnC8K-OJaA-WBx6-JLYz-HQU3-W8kb-0LHbXy   1          1          0.4%
      ```
    * #### Describe `LocalPV-LVM` volume
      ```bash
      $ kubectl openebs describe vol pvc-9999274f-ad01-48bc-9b21-7c51b47a870c
  
      pvc-9999274f-ad01-48bc-9b21-7c51b47a870c Details :
      ------------------
      Name            : pvc-9999274f-ad01-48bc-9b21-7c51b47a870c
      Namespace       : openebs
      AccessMode      : ReadWriteOnce
      CSIDriver       : local.csi.openebs.io
      Capacity        : 4Gi
      PVC             : csi-lvmpv
      VolumePhase     : Bound
      StorageClass    : openebs-lvmpv
      Version         : ci
      Status          : Ready
      VolumeGroup     : lvmvg
      Shared          : no
      ThinProvisioned : no
      NodeID          : worker-sh1
      ```