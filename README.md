# IRIS

<img src="https://github.com/olegsu/iris/raw/master/Iris.jpg" width="200">

`In Greek mythology, Iris is the personification of the rainbow and messenger of the gods.`


[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/olegs-codefresh/olegsu%2Firis%2Firis?branch=master&type=cf-2)]( https://g.codefresh.io/repositories/olegsu/iris/builds?filter=trigger:build;branch:master;service:5b69b8145904b871b671a6cf~iris)

Easily configure webhooks on Kubernets events using highly customize filters

* This project is not stable yet and may be changed anytime without any notice.

## Run in cluster
### Using Helm
* clone or fork this repository
* create your iris.yaml file
* install chart from local directory `helm install ./iris --values ./iris.yaml`
* by default the chart will be installed into namespace `iris`, see default values to overwrite it

## Build
### Locally
* Limitations:
  * Execute out of cluster `iris run --help`
  * Execute on non GCP cluster
* `make install`
* `make build`

Quick example:

In this example we will configure to listen on any Kubernetes event that been reported by the pod controller and matched to the filter will be sent to the destination.



```yaml
filters:
  - name: MatchDefaultNamespace
    type: namespace
    namespace: default
  - name: MatchPodKind
    type: jsonpath
    path: $.involvedObject.kind
    value: Pod

destinations:
  - name: prod
    url: http://localhost

integrations:
  - name: Report
    destination: 
    - prod
    filters:
    - MatchPodKind
    - MatchDefaultNamespace
```
On this example we will configure to listen on any Kubernetes event that been reported by the pod controller and matched to the filter will be sent to the destination.


## Filters
Set of rules that will be applied on each Kubernetes event
Kubernetes event that will pass all required filters will be passed to the destination to be reported
## Destinations
Each destination is an api endpoint where the Kubernetes event will be sent
## Integrations
Connecting between filters and destinations

## Execute Codefresh pipelines
With Iris, you can execute Codefresh pipelines.
Add destinations with codefresh type:
```yaml
  - name: ExecuteCodefreshPipeline
    type: Codefresh
    pipeline: PIPELINE_NAME
    cftoken: API_TOKEN
    branch: master
```