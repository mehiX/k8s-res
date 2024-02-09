# Kubernetes inspection tool (mainly resource)

## Install

```shell
go install github.com/mehix/kuberes@latest
kuberes version
```

## View resources declared in YAML manifest

Takes in input a series of YAML manifests, all in 1 file.

```shell
cat internal/aggr/testdata_confluentinc.yaml | kuberes declared
```

Test with a Helm repo:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm template bitnami/wordpress | kuberes declared --src k8s
```
