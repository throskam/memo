#!/bin/bash

set -e

mkdir -p .data/pgadmin
chmod 777 .data/pgadmin

mkdir -p .data/postgres
chmod 777 .data/postgres

systemd-run --scope --user -p "Delegate=yes" kind create cluster --name memo --config deployments/kind-config.yaml

sleep 15

podman update --pids-limit=16384 "$(podman ps --filter name=memo-control-plane --format '{{.ID}}')"

kubectx kind-memo

kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.yaml

sleep 30

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

kubectx -