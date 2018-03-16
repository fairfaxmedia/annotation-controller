# Annotation Controller

Controller that defines an ```annotation``` type in Kubernetes.

This allows the type to manage Annotations of other types within Kubernetes.


## Example

Defining an object of the type Annotation targeting the namespace object.

```
apiVersion: example.com/v1
kind: Annotation
metadata:
  name: example-annotation
  namespace: annotation-controller
spec:
  targets:
  - Data:
      kit: kat
    Kind: namespace
```

Will modify the annotations of the target resource

```
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    ping: pong
```

## Supported Types

|Type|Operations|
|------------------|-----------|
|`namespace`|CREATE/UPDATE|

## Building

```make```

## Running

```make deploy```

## Cleanup

```make delete```