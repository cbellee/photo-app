STORAGE_ACCOUNT_NAME=storjn45dot42kw3k
FOLDER_NAME='/mnt/c/Users/cbellee/Desktop/La_Moye_School_Photos'

# deelte existing images in all containers
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s thumbs
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s uploads
az storage blob delete-batch --account-name $STORAGE_ACCOUNT_NAME -s images

# uplod images
az storage blob upload-batch --account-name $STORAGE_ACCOUNT_NAME -d uploads -s $FOLDER_NAME --metadata collection=2022 album=nature --overwrite
