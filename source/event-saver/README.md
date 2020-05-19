# event-saver

The event saver is a simple demo application which subscribes to a pub/sub topic and saves all of the received events into storage. This section will overview the local development experience. 

## Prerequisites

To run this demo locally, you will have to install [Dapr](https://github.com).

## Setup

Start by cloning this repo

```shell
git clone git@github.com:mchmarny/app-portability-using-dapr.git
```

Navigate to the source repository 

```shell
cd source/event-saver
```

## Run

Launch the app using Dapr:

```shell
dapr run \
  --app-id event-saver \
  --app-port 8080 \
  --protocol http \
  --port 3500 \
  go run handler.go main.go
```

If everything goes well when launch these services using Dapr you will see following message:

```shell
ℹ️  Updating metadata for app command: dist/event-saver
✅  You're up and running! Both Dapr and your app logs will appear here.
```

## Demo 

Post some content unto the topic 

```shell
# print verbose since the publish post doesn't return any content
curl -v -H "Content-Type: application/json" \
    "http://localhost:3500/v1.0/publish/events" \
    -d '{ "id": "id1", "message": "hello" }'
```

> Notice that Dapr will output application logs along its own but in different color to help developer debug any issues

To query the local state store, first connect to Dapr Redis

```shell
docker run --rm -it --link dapr_redis redis redis-cli -h dapr_redis
```

Then query the store for keys 

```shell
KEYS event-saver*
```

Should return something similar to this

```shell
1) "event-saver||626cdb94-b163-46e6-9215-115dbc50ed58"
2) "event-saver||8d6ab10b-cf2a-47b0-8616-aea518361836"
```

If you want to view the saved content use one of the above keys 

```shell
HGET event-saver||626cdb94-b163-46e6-9215-115dbc50ed58 data
```

The saved content will look something like this

```shell
{ "id": "1", "message": "hello" }
```

## Next

* [Kubernates](../../deployment/k8s)
* [Azure Container Instances](../../deployment/aci)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../../LICENSE)



