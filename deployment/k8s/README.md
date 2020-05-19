# Kubernetes Deployment

This document will overview the sample app deployment into Kubernetes. For illustration purposes, all commands in this document will be based on Microsoft Azure but any modern Kubernates cluster will do.

## Install dapr

See instructions how to install dapr into your Kubernetes cluster [here](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md#installing-dapr-on-a-kubernetes-cluster)

## Deploy components for backing services 

From withing the [deployment/k8s](deployment/k8s) directory first deploy the components

```shell
kubectl apply -f component/store.yaml
kubectl apply -f component/topic.yaml
```

## Deploying demo 

```shell
kubectl apply -f .
```

You can check on the status of your deployment like this: 

```shell
kubectl get pods -l app=event-saver
```

The results should look similar to this (make sure each pod has READY status 2/2)

```shell
NAME                           READY   STATUS    RESTARTS   AGE
event-saver-78bb9584ff-ngcgp   2/2     Running   0          18s
```

Wait until the service has been configured with a public IP (has `EXTERNAL-IP`)

```shell
kubectl get svc -l app=event-saver -w
```

Than just capture the IP for ease of access 

```shell
export APP_IP=$(kubectl get svc -l app=event-saver -o jsonpath="{.items[0].status.loadBalancer.ingress[0].ip}")
```


If everything went OK, you should be able to access the user service thru Dapr API 

```shell
curl -H "Content-Type: application/json" \
    "http://${APP_IP}:3500/v1.0/invoke/event-saver/method/"
```

Response should look something like this (notice the localhost service invocation)

```json
{ "requestor": "127.0.0.1:8080", "version": "v0.1.5" }
```

> Note, the direct access to the user service is disabled.

To view logs from the Dapr container 

```shell
kubectl logs -l app=event-saver -c daprd
```

## Demo  

Post some content unto the topic 

```shell
# print verbose since the publish post doesn't return any content
curl -v -H "Content-Type: application/json" \
    "http://${APP_IP}:3500/v1.0/publish/events" \
    -d '{ "id": "1", "message": "hello from k8s" }'
```

The response should include success status code `HTTP/1.1 200 OK`


## Services

### Kafka

To view the Kafka topic, first, list the kafka-cat pod  

```shell
ubectl get pods -l app=kafka-cat
```

And then `exec` into it using the pod name

```shell
kubectl exec -ti kafka-cat-84b9cfdf45-glhtr bash
```

And connect to the topic  

```shell
kafkacat -b kafka-service:9092 -t events
```

### Redis 

To query the local state store, first connect pod

```shell
kubectl exec -ti redis-0 bash
```

Start Redis CLI 

```shell
redis-cli -h redis
```

Then query the store for keys 

```shell
KEYS event-saver*
```

Should return something similar to this

```shell
1) "event-saver||626cdb94-b163-46e6-9215-115dbc50ed58"
```

If you want to view the saved content use one of the above keys 

```shell
HGET event-saver||626cdb94-b163-46e6-9215-115dbc50ed58 data
```

The saved content will look something like this

```shell
{ "id": "1", "message": "hello" }
```

## Clean up 

First, delete the deployment and service 

```shell
kubectl delete -f .  --ignore-not-found
```

Then delete the components 

```shell
kubectl delete -f ../../component --ignore-not-found
```

## Next

* [Azure Container Instances](../../deployment/aci)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../../LICENSE) 

