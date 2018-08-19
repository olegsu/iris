# IRIS

<img src="https://github.com/olegsu/iris/raw/master/Iris.jpg" width="200" align="right">

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/olegs-codefresh/olegsu%2Firis%2Firis?branch=master&type=cf-2)]( https://g.codefresh.io/repositories/olegsu/iris/builds?filter=trigger:build;branch:master;service:5b69b8145904b871b671a6cf~iris)
[![Go Report Card](https://goreportcard.com/badge/github.com/olegsu/iris)](https://goreportcard.com/report/github.com/olegsu/iris)
[![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/kube-iris/Lobby)
[![codecov](https://codecov.io/gh/olegsu/iris/branch/master/graph/badge.svg)](https://codecov.io/gh/olegsu/iris)

<sub>**_In Greek mythology, Iris is the personification of the rainbow and messenger of the gods._**</sub>

- [IRIS](#iris)
    - [Run in cluster](#run-in-cluster)
        - [Using Helm](#using-helm)
    - [Build](#build)
    - [Filters](#filters)
        - [Reason](#reason)
        - [Namespace](#namespace)
        - [JSONPath](#jsonpath)
        - [Labels](#labels)
        - [Any](#any)
    - [Destinations](#destinations)
        - [Default](#default)
        - [Codefresh](#codefresh)
    - [Integrations](#integrations)

Easily configure webhooks on Kubernetes events using highly customizable filters

* This project is not stable yet and may be changed anytime without any notice.

## Run in cluster
### Using Helm
1. clone or fork this repository

```
$ git clone https://github.com/olegsu/iris.git
$ cd iris
```

2. create your iris.yaml file

```
$ cat << EOF > ./iris.yaml
filters:
  - name: MatchIrisNamespace
    type: namespace
    namespace: iris
  - name: MatchPodKind
    type: jsonpath
    path: $.involvedObject.kind
    value: Pod

destinations:
  - name: Webhook
    url: http://webhook-pod.iris.svc.cluster.local:8080

integrations:
  - name: Report
    destinations:
    - Webhook
    filters:
    - MatchPodKind
    - MatchIrisNamespace
EOF
```

3. install chart from local directory

```
$ helm install ./iris --values ./iris.yaml
```

by default the chart will be installed into namespace `iris`, see default values to overwrite it

4. If you want to webhook test, you apply under yaml.

```
$ #for exmaple...
$ cat << EOF > ./webhook-pod.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: webhook-pod
  namespace: iris
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webhook-pod
  template:
    metadata:
      labels:
        app: webhook-pod
    spec:
      containers:
      - name: webhook-pod
        image: nnao45/webhook-pod #Print that be received POST JSON contents with exposing port 8080
        ports:
        - name: http
          containerPort: 8080

---
apiVersion: v1
kind: Service
metadata:
  name: webhook-pod
  namespace: iris
spec:
  clusterIP: None
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: webhook-pod
EOF
$ kubectl apply -f webhook-pod.yaml
```

And `kubectl logs *<your pods>* -n iris`  
Print, Filtered events.

## Build
* clone or fork this repo
* `make install`
* `make build`
* Limitations:
  * Execute out of cluster `iris run --help`
  * Execute on non GCP cluster

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
    destinations: 
    - prod
    filters:
    - MatchPodKind
    - MatchDefaultNamespace
```

## Filters
Set of rules that will be applied on each [Kubernetes event](https://github.com/kubernetes/api/blob/master/core/v1/types.go#L4501).  
Kubernetes event that will pass all required filters will be passed to the destination to be reported  
Types of filters:
### Reason
Reason filter is a syntactic sugar for [JSONPath](#jsonpath) filter with `path: $.reason` and `value: {{reason}}`
```yaml
filters:
  - name: PodScheduled
    reason: "Scheduled"
```

### Namespace
Namespace filter is a syntactic sugar for [JSONPath](#jsonpath) filter with `path: $.metadata.namespace` and `value: {{reason}}`
```yaml
filters:
  - name: FromDefaultNamespace
    namespace: default
```

### JSONPath
With JSONPath gives the ability to match any field from [Kubernetes event](https://github.com/kubernetes/api/blob/master/core/v1/types.go#L4501).
The value from the fields can be matched to exec value using `value: {{value}}` or matched by regex using `regexp: {{regexp}}`
```yaml
filters:
  # Match to Warning event type
  - name: WarningLevel
    type: jsonpath
    path: $.type
    value: Warning
  # Match to any event that the name matched to regexp /tiller/
  - name: MatchToRegexpTiller
    type: jsonpath
    path: $.metadata.name
    regexp: tiller
```



### Labels
Labels filter will try to get the original resource from the event with the given filters.
The filter considers as passed if any resource were found
```yaml
filters:
   - name: MatchLabels
     type: labels
     labels:
       app: helm
```

### Any
```yaml
filters:
  - name: WarningLevel
    type: any
    filters:
    - FromDefaultNamespace
    - WarningLevel
```

## Destinations
Each destination is an api endpoint where the Kubernetes event will be sent
Types of destinations:
### Default
The default destinations will attempt to send POST request with the event json in the request body
If `secret` is given, hash string will be calculated using the given key on the request body and the result will be set in the request header as `X-IRIS-HMAC: hash`
```yaml
destinations:
  - name: Webhook
    url: https://webhook.site
    secret: SECRET
```
### Codefresh
With Iris, you can execute Codefresh pipelines.
Add destinations with Codefresh type:
* name: pipeline full name can be found easily using [Codefresh CLI](https://codefresh-io.github.io/cli/) - `codefresh get pipelines`
* branch: which branch of the repo should be cloned
* cftoken: Token to Codefresh API can be generated in [Account settings/Tokens](https://g.codefresh.io/account-conf/tokens) view
```yaml
  - name: ExecuteCodefreshPipeline
    type: Codefresh
    pipeline: PIPELINE_NAME
    cftoken: API_TOKEN
    branch: master
```

## Integrations
Connecting between filters and destinations
```yaml
integrations:
  - name: Report
    destinations:
    - {{name of destination}}
    filters:
    - {{name of filters to apply}}
```
