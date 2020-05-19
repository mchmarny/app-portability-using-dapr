# App Portability Using Dapr

Demo showing same apps being deployed across different runtime environments using [Dapr](https://dapr.io/) to bind app's backing state and pub/sub services.

> Dapr is an event-driven, portable runtime for building microservices. Whether developing locally, deploying to your own Kubernates cluster, or using a fully managed runtime environment like Azure Container Instances, Dapr provides a consistent HTTP and gRPC APIs for the most common building blocks of modern microservices: state, resource bindings, and pub/sub messaging. 

## Local Development

This demo will use a simple microservice application located in the [source/event-saver](source/event-saver) directory to subscribe to pub/sub topic and save all received events into state store. The instructions on how to setup your development environment and run this app locally are available [here](source/event-saver)

## Deployments

To show the ease of deploying the locally developed application to different runtime environments follow one of these links:

* [Kubernates](deployment/k8s)
* [Azure Container Instances](deployment/aci)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE) 


