# FinOps NATS Subscriber
This repository subscribes on the nats server topic specified in the environment variable `subTopic`. It then receives Azure optimizations from [Krateo FinOps HTTP Rest Queue](https://github.com/krateoplatformops/finops-http-rest-queue). The received optimizations are compiled into a Custom Resource managed by the [Krateo FinOps Operator VM Manager](https://github.com/krateoplatformops/finops-operator-vm-manager).

## Configuration
You must configure four parameters in the values.yaml file:
 - `subTopic`: the topic that the subscriber will receive data on;
 - `optNamespace`: the namespace where the "finops-operator-vm-manager" is deployed;
 - `optSecretName`: the name of the secret containing the token for the Azure REST API;
 - `optSecretNamespace`: the namespace of the secret containing the token for the Azure REST API;

## How to install
The installation can be performed using HELM:
```sh
$ helm repo add krateo https://charts.krateo.io
$ helm repo update krateo
$ helm install finops-nats-subscriber krateo/finops-nats-subscriber
```