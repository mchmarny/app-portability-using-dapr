# Azure App Service  Deployment 

This document will overview the sample app deployment into Azure App Service (AAS).

## Prerequisite

* [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)

Also, to simplify all the scripts in this doc, set a few `az` CLI defaults:

```shell
az account set --subscription <id or name>
az configure --defaults location=<prefered location> group=<your resource group>
```

## Backing Services 

The setup of the backing services is beyond the scope of this document but you can find detail instructions on how to configure Azure services for state and pub/sub below:

* [State store configure with Azure Table Storage](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal). 
* [PubSub topic configure using Azure Service Bus](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal).


## Setup

### App Service Plan

To deploy Dapr into AAS you will need to first create a app service plan 

```shell
export PLAN_NAME="eventsaver"
```

> assumes your resource group and location defaults are already set 

```shell
az appservice plan create --name $PLAN_NAME --sku S1 --is-linux
```

### File Share for Components 


The other thing you will need to deploy to AAS is a storage for Dapr components. 

> Storage account name needs to be unique, 3-24 chars long, and containing only alphanumerics. You may need to update the bellow `SANAME` variable 

```shell
export SA_NAME="eventsaveraas"
```

From withing the [deployment/aas](deployment/aas) directory, first, create a storage account

> assumes your resource group and location defaults are already set 

```shell
az storage account create --name $SA_NAME --sku Standard_LRS
```

Create a file share

```shell
az storage share create --name components --account-name $SA_NAME
```

Upload the Dapr component files

```shell
az storage file upload --account-name $SA_NAME --share-name components --source component/store.yaml
az storage file upload --account-name $SA_NAME --share-name components --source component/topic.yaml
```

Capture storage key

```shell
az storage account keys list --account-name $SA_NAME --query "[0].value" --output tsv
```

Now update `volumes[components].azureFile.storageAccountKey` in [deployment.yml](./deployment.yml) file so that ACI can mount it


## Deployment

Once the storage is set up, you can deploy the app as many time as you need

```shell
az webapp create --plan $PLAN_NAME \
                 --name event-saver \
                 --multicontainer-config-type compose \
                 --multicontainer-config-file deployment.yml
```

The result should look something like this 

```shell
az webapp list -o table
Name         Location    State    ResourceGroup    DefaultHostName                AppServicePlan
-----------  ----------  -------  ---------------  -----------------------------  ----------------
event-saver  West US 2   Running  dev_machmarn     event-saver.azurewebsites.net  eventsaver
```

### Storage 

> TODO: figure out the circular logic here. The volume mounted by AAC has be empty ot ACC will delete all things in there while Dapr needs this files ot be there when starting 

```shell
az webapp config storage-account add --name <app_name> --custom-id <custom_id> --storage-type AzureFiles --share-name <share_name> --account-name <storage_account_name> --access-key "<access_key>" --mount-path <mount_path_directory>
```

### Hostname 

Than just capture the IP for ease of access 

```shell
export APP_HOSTNAME=$(az webapp show -n event-saver --query "hostNames[0]" -o tsv)
```

If everything went OK, you should be able to access the user service thru Dapr API 

```shell
curl -H "Content-Type: application/json" \
    "http://${APP_HOSTNAME}:3500/v1.0/invoke/event-saver/method/"
```

Response should look something like this (notice the localhost service invocation)

```json
{ "requestor": "127.0.0.1:8080", "version": "v0.1.5" }
```

To view logs from the Dapr container 

```shell
az webapp log show --name event-saver
```

## Demo 

Post some content unto the topic 

```shell
# print verbose since the publish post doesn't return any content
curl -v -H "Content-Type: application/json" \
    "http://${APP_HOSTNAME}:3500/v1.0/publish/events" \
    -d '{ "id": "33", "message": "hi from aci" }'
```

The response should include success status code `HTTP/1.1 200 OK`


## Data

To view the data, first navigate to Service Bus topic in Azure Portal and see the message delivery count

![](../../image/topic.png)

Similarly, navigate to Storage Account section of the Azure Portal and see the event content saved into the table store 

![](../../image/state.png)

> for readability during demo, the JSON messages are saved in the state store as strings

That's it, I hope you found it helpful. 

## Clean up 

```shell
az webapp delete --name event-saver
```

## Next

## Next

* [Kubernates](../../deployment/k8s)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../../LICENSE) 

