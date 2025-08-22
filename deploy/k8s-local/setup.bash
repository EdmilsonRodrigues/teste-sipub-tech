#!/bin/bash

snap install microk8s --classic
cd helm || exit 1

microk8s.helm3 dependency build
microk8s.kubectl apply -f 'https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml'
microk8s.kubectl apply -f 'https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml'
microk8s.kubectl annotate storageclass local-path storageclass.kubernetes.io/is-default-class=true
microk8s.helm3 install movies-api .


