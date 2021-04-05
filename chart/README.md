# microservice-authentication Helm Charts

## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

## Get Repo Info

```console
helm repo add microservice-authentication <path to chartmuseum or ArtifactHub>
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart directly from the Github repo with the release name `microservice-authentication`:

```console
helm install microservice-authentication chart/ --values chart/values.yaml
```

To install the chart from ArtifactHub with the release name `microservice-authentication`:

```console
helm install microservice-authentication <org name>/microservice-authentication --version <version> 
```

## Uninstalling the Chart

To uninstall/delete the microservice-authentication deployment:

```console
helm delete microservice-authentication
```
