<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# JIVA Storage Engine Commands

## Table of Contents
* [Jiva](#jiva)
    * [Get Jiva volumes](#get-jiva-volumes)
    * [Describe Jiva volumes](#describe-jiva-volumes)

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
      $ kubectl openebs describe volume pvc-d6349231-a863-4d01-962e-24a4b372050d
        
      pvc-d6349231-a863-4d01-962e-24a4b372050d Details :
      -----------------
      NAME            : pvc-d6349231-a863-4d01-962e-24a4b372050d
      ACCESS MODE     : ReadWriteOnce
      CSI DRIVER      : jiva.csi.openebs.io
      STORAGE CLASS   : openebs-jiva-csi-sc
      VOLUME PHASE    : Bound
      VERSION         : 2.11.0
      JVP             : jivavolumepolicy
      SIZE            : 4.0GiB
      STATUS          : RW
      REPLICA COUNT	: 2
    
      Portal Details :
      ------------------
      IQN              :  iqn.2016-09.com.openebs.jiva:pvc-d6349231-a863-4d01-962e-24a4b372050d
      VOLUME NAME      :  pvc-d6349231-a863-4d01-962e-24a4b372050d
      TARGET NODE NAME :  node2-virtual-machine
      PORTAL           :  10.106.182.76:3260
    
      Replica Details :
      -----------------
      NAME                                                          STATUS   VOLUME                                     CAPACITY   STORAGECLASS       AGE
      openebs-pvc-d6349231-a863-4d01-962e-24a4b372050d-jiva-rep-0   Bound    pvc-19af6af5-34e1-4ab4-8401-07269c4084fa   4.0GiB     openebs-hostpath   6d5h
      openebs-pvc-d6349231-a863-4d01-962e-24a4b372050d-jiva-rep-1   Bound    pvc-50a75a4d-a809-45e0-8b48-167e8e5a2c17   4.0GiB     openebs-hostpath   6d5h
      ```
      