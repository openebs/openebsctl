<img width="300" align="right" alt="OpenEBS Logo" src="https://raw.githubusercontent.com/cncf/artwork/master/projects/openebs/stacked/color/openebs-stacked-color.png" xmlns="http://www.w3.org/1999/html">

# OpenEBSCTL


[![Go Report Card](https://goreportcard.com/badge/github.com/openebs/openebsctl?)](https://goreportcard.com/report/github.com/openebs/openebsctl)
[![Contributors](https://img.shields.io/github/contributors/openebs/openebsctl)](https://github.com/openebs/openebsctl/graphs/contributors)
[![release](https://img.shields.io/github/release-pre/openebs/openebsctl.svg)](https://github.com/openebs/openebsctl/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/mum4k/termdash/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/downloads/openebs/openebsctl/total.svg)](https://github.com//openebs/openebsctl/releases)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopenebs%2Fopenebsctl?ref=badge_shield)
[![codecov.io](https://codecov.io/github/openebs/openebsctl/coverage.svg?branch=develop)](https://codecov.io/github/openebs/openebsctl?branch=develop)


OpenEBSCTL is a kubectl plugin to manage OpenEBS storage components.


## Project Status

**Alpha**. Under active development and seeking [contributions from the community](#contributing).
The CLI currently supports managing `LocalPV-LVM`, `LocalPV-ZFS` and `LocalPV-HostPath` Engines.

## Table of Contents
* [Installation](#installation)
* [Build](#build)
* [Code Walkthrough](#code-walkthrough)
* [Usage](#usage)
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
   ```

## Build

- Clone this repo to your system. `git clone https://github.com/openebs/openebsctl`
- `cd openebsctl`
- Run `make openebsctl`

## Usage

* ```bash
  $ kubectl openebs
  kubectl openebs is a a kubectl plugin for interacting with OpenEBS storage components such as storage(zfspools, volumegroups), volumes, pvcs.
  Find out more about OpenEBS on https://openebs.io/
  
  Usage:
  openebs [command]
  
  Available Commands:
  cluster-info Show component version, status and running components for each installed engine
  completion   Outputs shell completion code for the specified shell (bash or zsh)
  describe     Provide detailed information about an OpenEBS resource
  get          Provides fetching operations related to a Volume/Storage
  help         Help about any command
  version      Shows openebs kubectl plugin's version
  
  Flags:
  -h, --help                help for openebs
  -c, --kubeconfig string   path to config file
  -v, --version             version for openebs
  
  Use "openebs [command] --help" for more information about a command.
  ```

* To know more about various engine specific commands check these:-
  * [LocalPV-LVM](docs/localpv-lvm/README.md)
  * [LocalPV-ZFS](docs/localpv-zfs/README.md)
  
## Contributing

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
