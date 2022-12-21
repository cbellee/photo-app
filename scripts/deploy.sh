#!/bin/bash

while getopts ":s" option; do
   case $option in
      s) skipBuild=1; # use '-s' cmdline flag to skip the container build step
   esac
done

LOCATION='australiaeast'
RG_NAME="photo-app-test-rg"
ENVIRONMENT=dev
SEMVER=0.1.0
REV=$(git rev-parse --short HEAD)
TAG="$ENVIRONMENT-$SEMVER-$REV"

DOMAIN_NAME="kainiindustries.net"
CLIENT_ID="beb06d91-7b0c-4214-8259-b798ed982d08"
PHOTO_APP_ID="api://a845082b-e22d-49a8-8abb-e8484609abd7"
UPLOAD_APP_ID="api://18911b98-3bf5-4a05-a417-8a12e496c9e5"
COLLECTION_API_ENDPOINT="https://https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/collections"
ALBUM_API_ENDPOINT="https://https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/albums"
PHOTO_API_ENDPOINT="https://https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/photos"
UPLOAD_API_ENDPOINT="https://https://photo.orangecoast-f46e9cc8.australiaeast.azurecontainerapps.io/uploads"
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

		cd ./scripts
fi

az deployment group create \
--resource-group $RG_NAME \
--name 'infra-deployment' \
--template-file ../infra/main.bicep \
--parameters ../infra/main.parameters.json \
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
--parameters grpcMaxRequestSizeMb=$GRPC_MAX_REQUEST_SIZE_MB

STORAGE_ACCOUNT_NAME=$(az deployment group show -g $RG_NAME -n 'infra-deployment' --query properties.outputs.storageAccount.value -o tsv)
PHOTO_URL=$(az deployment group show -g $RG_NAME -n 'infra-deployment' --query properties.outputs.photoUrl.value -o tsv)

# delete existing images in all containers
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s thumbs
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s uploads
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s images

# upload images and add metadata & index tags

for f in ../images/nature/*; 
do 
fileName=$(basename ${f})
az storage blob upload \
--account-name "$STORAGE_ACCOUNT_NAME" \
--file "$f" \
--container-name uploads \
--metadata collection=2022 album=nature \
--tags collection=2022 album=nature \
--name "2022/nature/$fileName"; 
done

for f in ../images/sport/*;
do 
fileName=$(basename ${f})
az storage blob upload \
--account-name "$STORAGE_ACCOUNT_NAME" \
--file "$f" \
--container-name uploads \
--metadata collection=2022 album=sport \
--tags collection=2022 album=sport \
--name "2022/sport/$fileName"; 
done

for f in ../images/soccer/*; 
do 
fileName=$(basename ${f})
az storage blob upload \
--account-name "$STORAGE_ACCOUNT_NAME" \
--file "$f" \
--container-name uploads \
--metadata collection=2022 album=soccer \
--tags collection=2022 album=soccer \
--name "2022/soccer/$fileName"; 
done

for f in ../images/school/*; 
do 
fileName=$(basename ${f})
az storage blob upload \
--account-name "$STORAGE_ACCOUNT_NAME" \
--file $f \
--container-name uploads \
--metadata collection=1979 album=school \
--tags collection=1979 album=school \
--name "1979/school/$fileName"; 
done


###############################

':
az storage blob upload-batch \
--account-name $STORAGE_ACCOUNT_NAME \
-d uploads \
-s ../images/nature \
--destination-path "2022/nature" \
--metadata collection=2022 album=nature \
--overwrite

az storage blob upload-batch \
--account-name $STORAGE_ACCOUNT_NAME \
-d uploads \
-s ../images/sport \
--destination-path "2022/sport" \
--metadata collection=2022 album=sport \
--overwrite

az storage blob upload-batch \
--account-name $STORAGE_ACCOUNT_NAME \
-d uploads \
-s ../images/soccer \
--destination-path "2022/soccer" \
--metadata collection=2022 album=soccer \
--overwrite

az storage blob upload-batch \
--account-name $STORAGE_ACCOUNT_NAME \
-d uploads \
-s ../images/school \
--destination-path "1979/school" \
--metadata collection=1979 album=school \
--overwrite
'

# /queryPhotos
sleep -s 30

':
# deploy React app to Azure storage
echo "patching react 'authConfig.js' file"

cat ../spa/src/authConfig_template.js | sed -e "s:{{CLIENT_ID}}:${CLIENT_ID}:g" \
-e "s+{{DOMAIN_NAME}}+${DOMAIN_NAME}+g" \
-e "s+{{REDIRECT_URI}}+photo-api.${DOMAIN_NAME}+g" \
-e "s+{{PHOTO_READ_SCOPE}}+${PHOTO_READ_SCOPE}+g" \
-e "s+{{PHOTO_WRITE_SCOPE}}+${PHOTO_WRITE_SCOPE}+g" \
-e "s+{{UPLOAD_READ_SCOPE}}+${UPLOAD_READ_SCOPE}+g" \
-e "s+{{UPLOAD_WRITE_SCOPE}}+${UPLOAD_WRITE_SCOPE}+g" \
-e "s+{{UPLOAD_API_ENDPOINT}}+${UPLOAD_API_ENDPOINT}+g" \
-e "s+{{PHOTO_API_ENDPOINT}}+${PHOTO_API_ENDPOINT}+g"  > ../spa/src/authConfig.js
'

curl "https://$PHOTO_URL/photos?album=sport&collection=2022"  | jq
