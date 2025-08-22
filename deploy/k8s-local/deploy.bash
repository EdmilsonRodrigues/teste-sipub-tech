#!/bin/bash

HOME=/home/ubuntu
VM_CPUS=2
VM_MEMORY=4GB
VM_IMAGE=ubuntu:22.04


VM_NAME=$1

if [[ $VM_NAME == "" ]]; then
    echo "missing vm name."
    exit 1
fi

lxc launch $VM_IMAGE $VM_NAME --vm -c limits.cpu=$VM_CPUS -c limits.memory=$VM_MEMORY

sleep 20


./update.bash $VM_NAME $HOME
lxc exec $VM_NAME -- bash -c "snap install microk8s --classic"
lxc exec $VM_NAME -- bash -c "cd ${HOME}/helm; microk8s.helm3 dependency build"
lxc exec $VM_NAME -- bash -c "microk8s.kubectl apply -f 'https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml'"
lxc exec $VM_NAME -- bash -c "microk8s.kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml"
lxc exec $VM_NAME -- bash -c "microk8s.kubectl annotate storageclass local-path storageclass.kubernetes.io/is-default-class=true"
lxc exec $VM_NAME -- bash -c "cd ${HOME}/helm; microk8s.helm3 install movies-api ."
