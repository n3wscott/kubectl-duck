# kubectl-duck
Ducktyping support for kubectl.

[![GoDoc](https://godoc.org/github.com/n3wscott/kubectl-duck?status.svg)](https://godoc.org/github.com/n3wscott/kubectl-duck)
[![Go Report Card](https://goreportcard.com/badge/n3wscott/kubectl-duck)](https://goreportcard.com/report/n3wscott/kubectl-duck)

_Work in progress._

## Installation

Use [krew](https://sigs.k8s.io/krew) plugin manager to install,

```shell script
kubectl krew install duck
kubectl duck --help
```

Or manually,

`kubectl-duck` can be installed via:

```shell
go get github.com/n3wscott/kubectl-duck/cmd/kubectl-duck
```

To update your installation:

```shell
go get -u github.com/n3wscott/kubectl-duck/cmd/kubectl-duck
```

## Usage

`kubectl-duck` has two command at the moment: `list` and `get`

### List 

```shell
List CustomResourceDefinitions that implement ducktypes.

Usage:
  kubectl-duck list [flags]

Examples:

  To list the known ducktypes:
  $ kubectl duck list

  To list which resources a given ducktype maps to:
  $ kubectl duck list <ducktype>


Flags:
  -h, --help   help for list

Global Flags:
  -v, --v Level   number for the log level verbosity
```

#### Examples
 
 ```shell
 $ kubectl duck list
 NAME                                SHORT NAME    SELECTOR
 duck.knative.dev/source             source        duck.knative.dev/source=true
 messaging.knative.dev/subscribable  subscribable  messaging.knative.dev/subscribable=true
 duck.knative.dev/addressable        addressable   duck.knative.dev/addressable=true
 ```
 
 ```shell
 $ kubectl duck list addressable
 NAME                                    KIND             DUCK         CREATED AT
 brokers.eventing.knative.dev            Broker           addressable  66d
 channels.messaging.knative.dev          Channel          addressable  66d
 inmemorychannels.messaging.knative.dev  InMemoryChannel  addressable  66d
 parallels.flows.knative.dev             Parallel         addressable  60d
 parallels.messaging.knative.dev         Parallel         addressable  66d
 routes.serving.knative.dev              Route            addressable  66d
 sequences.flows.knative.dev             Sequence         addressable  60d
 sequences.messaging.knative.dev         Sequence         addressable  66d
 services.serving.knative.dev            Service          addressable  66d
 tasks.n3wscott.com                      Task             addressable  62d
 ```

### Get

```shell
Get the resource instances related to a ducktype.

Usage:
  kubectl-duck get [flags]

Examples:

  To get resource instances that are of the given ducktype:
  $ kubectl duck get [ducktype]


Flags:
  -h, --help   help for get

Global Flags:
  -v, --v Level   number for the log level verbosity
```

#### Examples
 
 ```shell
$ kubectl duck get addressable
  Broker brokers.eventing.knative.dev/v1alpha1
  NAMESPACE  NAME            READY  REASON  AGE
  default    Broker/default  True           62d
  
  
  InMemoryChannel inmemorychannels.messaging.knative.dev/v1alpha1
  NAMESPACE  NAME                                 READY  REASON  AGE
  default    InMemoryChannel/default-kne-ingress  True           60d
  default    InMemoryChannel/default-kne-trigger  True           60d
  
  
  
  Route routes.serving.knative.dev/v1
  NAMESPACE  NAME                  READY  REASON  AGE
  default    Route/classifier      True           60d
  default    Route/even            True           60d
  default    Route/odd             True           60d
  
  
  Service services.serving.knative.dev/v1
  NAMESPACE  NAME                READY  REASON  AGE
  default    Service/classifier  True           60d
  default    Service/even        True           60d
  default    Service/odd         True           60d
 ```

## TODO:

- [] Add support for dynamic ducktype lookup.
- [] Add support for known ducktypes local file.
 
## Authors
 
Scott Nichols [@n3wscott](https://twitter.com/n3wscott).
 
## License
 
Apache 2.0. See [LICENSE](./LICENSE).
