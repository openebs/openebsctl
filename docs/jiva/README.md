<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# JIVA Storage Engine Commands

## Table of Contents
* [Jiva](#jiva)
  * [Get Jiva volumes](#get-jiva-volumes)
  * [Describe Jiva volumes](#describe-jiva-volumes)
  * [Describe Jiva PVCs](#describe-jiva-pvcs)

* #### `Jiva`
    * #### Get `Jiva` volumes
      ```bash
      $ kubectl openebs get volumes --cas-type=jiva
      NAMESPACE   NAME                                       STATUS   VERSION   CAPACITY   STORAGE CLASS         ATTACHED   ACCESS MODE     ATTACHED NODE
      openebs     pvc-478a8329-f02d-47e5-8288-0c28b582be25   RW       2.9.0     4Gi        openebs-jiva-csi-sc   Released   ReadWriteOnce   minikube-2
      ```
      Note: For volumes not attached to any application, the `ATTACH NODE` would be shown as `N/A`.

    * #### Describe `Jiva` volumes
      ```bash
      $ kubectl openebs describe volume pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
        
      pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca Details :
      -----------------
      NAME            : pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      ACCESS MODE     : ReadWriteOnce
      CSI DRIVER      : jiva.csi.openebs.io
      STORAGE CLASS   : openebs-jiva-csi-sc
      VOLUME PHASE    : Bound
      VERSION         : 2.12.1
      JVP             : jivavolumepolicy
      SIZE            : 4.0GiB
      STATUS          : RW
      REPLICA COUNT	: 1
    
      Portal Details :
      ------------------
      IQN              :  iqn.2016-09.com.openebs.jiva:pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      VOLUME NAME      :  pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      TARGET NODE NAME :  minikube
      PORTAL           :  10.108.189.51:3260
    
      Controller and Replica Pod Details :
      -----------------------------------
      NAMESPACE   NAME                                                              MODE   NODE       STATUS    IP            READY   AGE
      jiva        pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-ctrl-64c964bvtbk5   RW     minikube   Running   172.17.0.9    1/1     8h25m
      jiva        pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-rep-0               RW     minikube   Running   172.17.0.10   1/1     8h25m
    
      Replica Data Volume Details :
      -----------------------------
      NAME                                                          STATUS   VOLUME                                     CAPACITY   STORAGECLASS       AGE
      openebs-pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-rep-0   Bound    pvc-009a193e-aa44-44d8-8b13-58859ffa734d   4.0GiB     openebs-hostpath   8h25m
      ```
      
    * #### Describe `Jiva` PVCs
      ```bash
      jiva-csi-pvc Details  :
      -------------------
      NAME               : jiva-csi-pvc
      NAMESPACE          : default
      CAS TYPE           : jiva
      BOUND VOLUME       : pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      ATTACHED TO NODE   : minikube
      JIVA VOLUME POLICY : jivavolumepolicy
      STORAGE CLASS      : openebs-jiva-csi-sc
      SIZE               : 4Gi
      JV STATUS          : RW
      PV STATUS          : Bound
    
      Portal Details :
      ------------------
      IQN              :  iqn.2016-09.com.openebs.jiva:pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      VOLUME NAME      :  pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca
      TARGET NODE NAME :  minikube
      PORTAL           :  10.108.189.51:3260
    
      Controller and Replica Pod Details :
      -----------------------------------
      NAMESPACE   NAME                                                              MODE   NODE       STATUS    IP            READY   AGE
      jiva        pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-ctrl-64c964bvtbk5   RW     minikube   Running   172.17.0.9    1/1     8h24m
      jiva        pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-rep-0               RW     minikube   Running   172.17.0.10   1/1     8h24m
    
      Replica Data Volume Details :
      -----------------------------
      NAME                                                          STATUS   VOLUME                                     CAPACITY   STORAGECLASS       AGE
      openebs-pvc-e974f45d-8b8f-4939-954a-607f60a8a5ca-jiva-rep-0   Bound    pvc-009a193e-aa44-44d8-8b13-58859ffa734d   4.0GiB     openebs-hostpath   8h24m
      ```