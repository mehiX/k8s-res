# Overview of resources declared in Kubernetes deployments

## Run

```shell
make
cat internal/aggr/testdata_confluentinc.yaml | go run ./main.go
```

Test with a Helm repo:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm template bitnami/wordpress | go run ./main.go --src k8s
```

## Install and use binary

```shell
go install github.com/mehix/k8s-resources
k8s-resources -h
```