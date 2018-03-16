#!/usr/bin/env bash

kubectl -n annoataion-controller delete annotation example-annotation
kubectl -n annotation-controller delete serviceaccount default
kubectl delete clusterrolebinding annotation-controller
kubectl delete crd annotations.example.com
kubectl delete ns annotation-controller
