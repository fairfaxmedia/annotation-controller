#!/usr/bin/env bash

kubectl create ns annotation-controller
kubectl apply -f artifacts/
kubectl -n annotation-controller run controller --image=quay.io/pickledrick/annotation-controller --replicas=1