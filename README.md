# Log-Cache Operator

Based on operator-framework-sdk Log-Cache Operator operates a Log-Cache cluster on top of Kubernetes API.

Create TLS assets

```sh
docker run -v "$PWD/output:/output" loggregator/certs
kubectl create secret generic logcachenozzle-tls --from-file=./output
kubectl create secret generic logcachescheduler-tls --from-file=./output
```

Deploy the Operator

```sh
kubectl create -f ./deploy
```

Create a new Log Cache cluster
```sh
kubectl create -f ./deploy/example
```

Make sure things start up ok
```sh
kubectl get pods
```

Make some changes
```sh
kubectl edit logcache example
```

Delete the cluster
```sh
kubectl delete logcache example
```
