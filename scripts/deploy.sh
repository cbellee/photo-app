while getopts ":s" option; do
   case $option in
      s) skipBuild=1; # skip container build step
   esac
done

RG_NAME='photo-app-rg'
LOCATION='australiaeast'
ENVIRONMENT=dev
SEMVER=0.1.0
REV=$(git rev-parse --short HEAD)
TAG="$ENVIRONMENT-$SEMVER-$REV"

RESIZE_API_NAME="resize"
RESIZE_API_PORT="443"
RESIZE_LOCAL_API_PORT="8443"

STORE_API_NAME="store"
STORE_API_PORT="443"
STORE_LOCAL_API_PORT="9443"

PHOTO_API_NAME="photo"
PHOTO_API_PORT="443"
PHOTO_LOCAL_API_PORT="10443"

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

MAX_THUMB_HEIGHT='128'
MAX_THUMB_WIDTH='96'
MAX_IMAGE_HEIGHT='640'
MAX_IMAGE_WIDTH='480'

echo "TAG: $TAG"

az group create --location $LOCATION --name $RG_NAME

if [[ $skipBuild != 1 ]]; then
	az deployment group create \
		--resource-group $RG_NAME \
		--name 'acr-deployment' \
		--parameters anonymousPullEnabled=true \
		--template-file ./infra/modules/acr.bicep
fi

ACR_NAME=$(az deployment group show --resource-group $RG_NAME --name 'acr-deployment' --query properties.outputs.acrName.value -o tsv)

if [[ $skipBuild != 1 ]]; then
	# build image in ACR
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
fi

az deployment group create \
	--resource-group $RG_NAME \
	--name 'infra-deployment' \
	--template-file ./infra/main.bicep \
	--parameters ./infra/main.parameters.json \
	--parameters location=$LOCATION \
	--parameters resizeApiName=$RESIZE_API_NAME \
	--parameters resizeApiPort=$RESIZE_API_PORT \
	--parameters photoApiName=$PHOTO_API_NAME \
	--parameters photoApiPort=$PHOTO_API_PORT \
	--parameters storeApiName=$STORE_API_NAME \
	--parameters storeApiPort=$STORE_API_PORT \
	--parameters tag=$TAG \
	--parameters acrName=$ACR_NAME \
	--parameters uploadsStorageQueueName=$UPLOADS_QUEUE_NAME \
	--parameters imagesStorageQueueName=$IMAGES_QUEUE_NAME \
	--parameters thumbsContainerName=$THUMBS_CONTAINER_NAME \
	--parameters imagesContainerName=$IMAGES_CONTAINER_NAME \
	--parameters uploadsContainerName=$UPLOADS_CONTAINER_NAME \
	--parameters maxThumbHeight=$MAX_THUMB_HEIGHT \
	--parameters maxThumbWidth=$MAX_THUMB_WIDTH \
	--parameters maxImageHeight=$MAX_IMAGE_HEIGHT \
	--parameters maxImageWidth=$MAX_IMAGE_WIDTH \
	--parameters databaseName=$COSMOSDB_NAME \
	--parameters containerName=$COSMOSDB_CONTAINER_NAME \
	--parameters partitionKey=$COSMOSDB_PARTITION_KEY \
	--parameters isExternalIngressEnabled='false'
