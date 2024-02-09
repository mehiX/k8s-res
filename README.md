# Kubernetes inspection tool (mainly resource)

## Install

```shell
go install github.com/mehix/kres@latest
kres version
```

## View resources declared in YAML manifest

Takes in input a series of YAML manifests, all in 1 file.

```shell
cat internal/aggr/testdata_confluentinc.yaml | kres declared
```

Test with a Helm repo:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm template bitnami/wordpress | kres declared --src k8s
```
