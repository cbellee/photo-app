#!/bin/bash

while getopts ":s" option; do
   case $option in
      s) skipBuild=1; # use '-s' cmdline flag to skip the container build step
   esac
done

LOCATION='australiaeast'
RG_NAME="photo-app-test-rg"
SEMVER='0.1.0'
REV=$(git rev-parse --short HEAD)
TAG="$SEMVER-$REV"

DOMAIN_NAME="kainiindustries.net"
CLIENT_ID="beb06d91-7b0c-4214-8259-b798ed982d08"
PHOTO_APP_ID="api://a845082b-e22d-49a8-8abb-e8484609abd7"
UPLOAD_APP_ID="api://18911b98-3bf5-4a05-a417-8a12e496c9e5"
PHOTO_READ_SCOPE="${PHOTO_APP_ID}/Photo.Read"
PHOTO_WRITE_SCOPE="${PHOTO_APP_ID}/Photo.Write"
UPLOAD_READ_SCOPE="${UPLOAD_APP_ID}/Upload.Read"
UPLOAD_WRITE_SCOPE="${UPLOAD_APP_ID}/Upload.Write"

RESIZE_API_NAME="resize"
RESIZE_API_PORT="443"
STORE_API_NAME="store"
STORE_API_PORT="443"
PHOTO_API_NAME="photo"
PHOTO_API_PORT="443"
GRPC_MAX_REQUEST_SIZE_MB="30"

RESIZE_API_IMAGE="$RESIZE_API_NAME:$TAG"
STORE_API_IMAGE="$STORE_API_NAME:$TAG"
PHOTO_API_IMAGE="$PHOTO_API_NAME:$TAG"

UPLOADS_QUEUE_NAME='uploads'
IMAGES_QUEUE_NAME='images'
THUMBS_CONTAINER_NAME='thumbs'
IMAGES_CONTAINER_NAME='images'
UPLOADS_CONTAINER_NAME='uploads'

COSMOSDB_NAME='photodb'
COSMOSDB_CONTAINER_NAME='photos'
COSMOSDB_PARTITION_KEY='/partitionKey'

MAX_THUMB_HEIGHT='96'
MAX_THUMB_WIDTH='128'
MAX_IMAGE_HEIGHT='480'
MAX_IMAGE_WIDTH='640'

az group create --location $LOCATION --name $RG_NAME

if [[ $skipBuild != 1 ]]; then
	az deployment group create \
		--resource-group $RG_NAME \
		--name 'acr-deployment' \
		--parameters anonymousPullEnabled=true \
		--template-file ../infra/modules/acr.bicep
fi

ACR_NAME=$(az deployment group show --resource-group $RG_NAME --name 'acr-deployment' --query properties.outputs.acrName.value -o tsv)

if [[ $skipBuild != 1 ]]; then

	cd ..
	echo "TAG: $TAG"

	# build image in ACR
	az acr login -n $ACR_NAME 

	docker build -t "$ACR_NAME.azurecr.io/$RESIZE_API_IMAGE" \
	--build-arg SERVICE_NAME=$RESIZE_API_NAME \
	--build-arg SERVICE_PORT=$RESIZE_API_PORT \
	-f ./Dockerfile .

	docker push "$ACR_NAME.azurecr.io/$RESIZE_API_IMAGE"

':
	az acr build -r $ACR_NAME \
		-t $RESIZE_API_IMAGE \
		--build-arg SERVICE_NAME=$RESIZE_API_NAME \
		--build-arg SERVICE_PORT=$RESIZE_API_PORT \
		-f ./Dockerfile .

	az acr build -r $ACR_NAME \
		-t $STORE_API_IMAGE \
		--build-arg SERVICE_NAME=$STORE_API_NAME \
		--build-arg SERVICE_PORT=$STORE_API_PORT \
		-f ./Dockerfile .

	az acr build -r $ACR_NAME \
		-t $PHOTO_API_IMAGE \
		--build-arg SERVICE_NAME=$PHOTO_API_NAME \
		--build-arg SERVICE_PORT=$PHOTO_API_PORT \
		-f ./Dockerfile .
'

		cd ../scripts
fi

az deployment group create \
--resource-group $RG_NAME \
--name 'infra-deployment' \
--template-file ../infra/main.bicep \
--parameters ../infra/main.parameters.json \
--parameters location=$LOCATION \
--parameters resizeApiName=$RESIZE_API_NAME \
--parameters resizeApiPort=$RESIZE_API_PORT \
--parameters tag=$TAG \
--parameters acrName=$ACR_NAME \
--parameters uploadsStorageQueueName=$UPLOADS_QUEUE_NAME \
--parameters thumbsContainerName=$THUMBS_CONTAINER_NAME \
--parameters imagesContainerName=$IMAGES_CONTAINER_NAME \
--parameters uploadsContainerName=$UPLOADS_CONTAINER_NAME \
--parameters maxThumbHeight=$MAX_THUMB_HEIGHT \
--parameters maxThumbWidth=$MAX_THUMB_WIDTH \
--parameters maxImageHeight=$MAX_IMAGE_HEIGHT \
--parameters maxImageWidth=$MAX_IMAGE_WIDTH \
--parameters grpcMaxRequestSizeMb=$GRPC_MAX_REQUEST_SIZE_MB
