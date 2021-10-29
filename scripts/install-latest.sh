#!/usr/bin/env bash

# Copyright 2020 The OpenEBS Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This script downloads and copies the latest openebsctl binary from github releases to /usr/local/bin
set -e

# Determine the arch/os combos before running the script

UNAME=$(uname)
ARCH=$(uname -m)

if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    echo "Sorry, this script is not supported for this OS."
    exit 1
fi

if [ "${ARCH}" = "i386" ] ; then
    XC_ARCH='x86_32'
elif [ "${ARCH}" = "x86_64" ] ; then
    XC_ARCH='x86_64'
elif [ "${ARCH}" = "aarch64" ] ; then
    XC_ARCH='arm64'
else
    echo "Sorry, this script is not supported for this arch: $XC_ARCH"
    exit 1
fi

# bool function to test if the user is root or not (POSIX only)
is_user_root () { [ "$(id -u)" -eq 0 ]; }

echo -e "\n\nGetting Latest Release ----->"

# GET Request to github API to fetch latest release tag
LATEST_TAG=$(curl --silent "https://api.github.com/repos/openebs/openebsctl/releases/latest"|grep '"tag_name":'|sed -E 's/.*"([^"]+)".*/\1/')

# Release Download link prefix
RELEASE_DOWNLOAD_LINK="https://github.com/openebs/openebsctl/releases/download/"

# Appending download link and latest tag for latest release link
LATEST_RELEASE_DOWNLOAD_LINK="$RELEASE_DOWNLOAD_LINK$LATEST_TAG"

# Appending binary name with latest release download for latest binary release link
LATEST_BINARY_DONWLOAD_LINK="$LATEST_RELEASE_DOWNLOAD_LINK/kubectl-openebs_""$LATEST_TAG""_"$UNAME"_$XC_ARCH.tar.gz"

cd /tmp

echo -e "\n\nDownloading Latest Release ----->\n\n"

wget -O openebsctl.tar.gz $LATEST_BINARY_DONWLOAD_LINK

tar -xvf openebsctl.tar.gz

echo -e "\n\nExtracted Latest Release ----->"

if is_user_root; then
    cp kubectl-openebs /usr/local/bin
else
    sudo cp kubectl-openebs /usr/local/bin
fi

echo -e "\n\nCopied Latest Release to usr/local/bin ----->"

echo -e "\n\nCleaning things ----->"

rm openebsctl.tar.gz && rm LICENSE && rm kubectl-openebs && rm README.md

echo -e "\n\nDone"
