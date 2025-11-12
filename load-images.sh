#!/bin/bash
set -e
CLUSTER=ticketing-demo

docker build -t api-gateway:latest ./api-gateway
docker build -t ticket-service:latest ./ticket-service
docker build -t waiting-room-service:latest ./waiting-room-service
docker build -t dashboard-service:latest ./dashboard-service
docker build -t website-dashboard:latest ./website-dashboard

k3d image import api-gateway:latest ticket-service:latest waiting-room-service:latest dashboard-service:latest website-dashboard:latest --cluster $CLUSTER
