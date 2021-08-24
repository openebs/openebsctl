<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# OpenEBSCTL


[![Go Report Card](https://goreportcard.com/badge/github.com/openebs/openebsctl?)](https://goreportcard.com/report/github.com/openebs/openebsctl)
[![Contributors](https://img.shields.io/github/contributors/openebs/openebsctl)](https://github.com/openebs/openebsctl/graphs/contributors)
[![release](https://img.shields.io/github/release-pre/openebs/openebsctl.svg)](https://github.com/openebs/openebsctl/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/mum4k/termdash/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/downloads/openebs/openebsctl/total.svg)](https://github.com//openebs/openebsctl/releases)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl?ref=badge_shield)



OpenEBSCTL is a kubectl plugin to manage OpenEBS storage components.


## Project Status

**Alpha**. Under active development and seeking [contributions from the community](#contributing).
The CLI currently supports managing `cStor`, `Jiva`, `LocalPV-LVM`, `LocalPV-ZFS` Cas-Engines.

## Table of Contents
* [Installation](#installation)
* [Build](#build)
* [Flags](#flags)
* [Usage](#usage)
  * [cStor](docs/cstor/README.md#cstor)
    * [Get cStor volumes](docs/cstor/README.md#get-cstor-volumes)
    * [Get cStor pools](docs/cstor/README.md#get-cstor-pools)
    * [Describe cStor volumes](docs/cstor/README.md#describe-cstor-volumes)
    * [Describe cStor pool](docs/cstor/README.md#describe-cstor-pool)
    * [Describe cStor PVCs](docs/cstor/README.md#describe-pvcs)
    * [Debugging cStor Volumes](docs/cstor/README.md#debugging-cstor-volumes)
  * [Jiva](docs/jiva/README.md#jiva)
    * [Get Jiva volumes](docs/jiva/README.md#get-jiva-volumes)
    * [Describe Jiva volumes](docs/jiva/README.md#describe-jiva-volumes)
    * [Describe Jiva PVCs](docs/jiva/README.md#describe-jiva-pvcs)
  * [LocalPV-LVM](docs/localpv-lvm/README.md#localpv-lvm)
    * [Get LocalPV-LVM volumes](docs/localpv-lvm/README.md#get-localpv-lvm-volumes)
    * [Get LocalPV-LVM VolumeGroups](docs/localpv-lvm/README.md#get-localpv-lvm-volumegroups)
    * [Describe LocalPV-LVM volumeGroups](docs/localpv-lvm/README.md#describe-localpv-lvm-volumeGroups)
    * [Describe LocalPV-LVM volumes](docs/localpv-lvm/README.md#describe-localpv-lvm-volumes)
  * [LocalPV-ZFS](docs/localpv-zfs/README.md#localpv-zfs)
    * [Get LocalPV-ZFS volumes](docs/localpv-zfs/README.md#get-localpv-zfs-volumes)
    * [Get LocalPV-ZFS Pools](docs/localpv-zfs/README.md#get-localpv-zfs-pools)
    * [Describe LocalPV-ZFS volumes](docs/localpv-zfs/README.md#describe-localpv-zfs-volumes)
    * [Describe LocalPV-ZFS pools](docs/localpv-zfs/README.md#describe-localpv-zfs-pools)
  * [BlockDevice](docs/cstor/README.md#blockdevice)
    * [Get BlockDevices by Nodes](docs/cstor/README.md#get-blockdevices-by-nodes)
* [Contributing](#contributing)

## Installation

OpenEBSCTL is available on Linux, macOS and Windows platforms.

* (**Recommended**) The latest binary can be installed via `krew`
  ```bash
  $ kubectl krew install openebs
  ...
  ...
  $ kubectl krew list
  PLUGIN    VERSION
  openebs    v0.2.0
  ...
  ...
  # to update the openebs plugin
  $ kubectl krew upgrade openebs
  ...
  ...
  ```

* Binaries for Linux, Mac and Windows are available as tarballs and zip in the [release](https://github.com/openebs/openebsctl/releases) page.
* Or, if you don't want to setup krew, you run the following to get latest version :-
   ```shell
   wget https://raw.githubusercontent.com/openebs/openebsctl/develop/scripts/install-latest.sh -O - | bash

## Build

- Clone this repo to your system. `git clone https://github.com/openebs/openebsctl`
- `cd openebsctl`
- Run `make openebsctl`
- Run `kubectl openebs [get|describe] [resource]` to use the plugin

### Flags

* `--openebs-namespace` :- to override the determination of `namespace` where storage engine is installed with the provided value.
* `--namespace, -n` :- to pass the namespace, if the resource is namespaced, like `pvc` etc.
* `--cas-type` :- to pass the cas-type, like cstor, jiva, localpv-lvm, localpv-zfs.

### Contributing

OpenEBS welcomes your feedback and contributions in any form possible.

- [Join OpenEBS community on Kubernetes Slack](https://kubernetes.slack.com)
    - Already signed up? Head to our discussions at [#openebs](https://kubernetes.slack.com/messages/openebs/)
- Want to raise an issue or help with fixes and features?
    - See [open issues](https://github.com/openebs/openebs/issues)
    - See [contributing guide](./CONTRIBUTING.md)
    - See [Project Roadmap](https://github.com/openebs/openebsctl/projects/1)
    - Checkout our existing [adopters](https://github.com/openebs/openebs/tree/master/adopters) and their [feedbacks](https://github.com/openebs/openebs/issues/2719).
    - Want to join our contributor community meetings, [check this out](https://hackmd.io/mfG78r7MS86oMx8oyaV8Iw?view).
- Join our OpenEBS CNCF Mailing lists
    - For OpenEBS project updates, subscribe to [OpenEBS Announcements](https://lists.cncf.io/g/cncf-openebs-announcements)
    - For interacting with other OpenEBS users, subscribe to [OpenEBS Users](https://lists.cncf.io/g/cncf-openebs-users)


For more details checkout [CONTRIBUTING.md](./CONTRIBUTING.md).



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl?ref=badge_large)
