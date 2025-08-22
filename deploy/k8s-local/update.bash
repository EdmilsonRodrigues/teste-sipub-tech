#!/bin/bash

VM_NAME=$1
HOME=$2
HELM_PATH=./helm
COMPRESSED_PATH=/tmp/helm.tar.xz
DESTINATION_PATH=$HOME/helm.tar.xz

tar -cvJf $COMPRESSED_PATH $HELM_PATH
lxc file push $COMPRESSED_PATH "${VM_NAME}${DESTINATION_PATH}"
lxc exec $VM_NAME -- bash -c "tar -xJf $DESTINATION_PATH -C $HOME"
lxc exec $VM_NAME -- bash -c "rm $DESTINATION_PATH"

