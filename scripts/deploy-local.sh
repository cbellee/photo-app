export UPLOADS_QUEUE_BINDING='uploadsQueue'
export IMAGES_QUEUE_BINDING='imagesQueue'
export COSMOS_BINDING='cosmosState'

export THUMBS_CONTAINER_STATE='thumbsState'
export IMAGES_CONTAINER_STATE='imagesState'
export UPLOADS_CONTAINER_STATE='uploadsState'

export IMAGES_CONTAINER_BINDING='imagesBinding'
export THUMBS_CONTAINER_BINDING='thumbsBinding'
export UPLOADS_CONTAINER_BINDING='uploadsBinding'

export COSMOS_DB="photodb"
export COSMOS_CONTAINER="photos"

export MAX_THUMB_HEIGHT=96
export MAX_THUMB_WIDTH=128
export MAX_IMAGE_HEIGHT=480
export MAX_IMAGE_WIDTH=640

source ./.env

cd ../cmd
cd ./resize
go build -o ../../bin
cd ..
cd ./store
go build -o ../../bin
cd ..
cd ./photo
go build -o ../../bin

cd ../..

dapr stop resize
dapr stop store
dapr stop photo

export SERVICE_PORT="8080"
export SERVICE_NAME="resize"
dapr run --app-id $SERVICE_NAME --app-protocol grpc --app-port $SERVICE_PORT ./bin/resize --components-path ./components --log-level debug &

export SERVICE_PORT="8081"
export SERVICE_NAME="store"
dapr run --app-id $SERVICE_NAME --app-protocol grpc --app-port $SERVICE_PORT ./bin/store --components-path ./components --log-level debug &

export SERVICE_PORT="8082"
export SERVICE_NAME="photo"
dapr run --app-id $SERVICE_NAME --app-protocol grpc --app-port $SERVICE_PORT --dapr-http-port 3000 ./bin/photo --components-path ./components --log-level debug &
