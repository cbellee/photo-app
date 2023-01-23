import { apiConfig, msalConfig } from "./authConfig";
import { InteractiveBrowserCredential } from "@azure/identity"
import { BlockBlobClient } from "@azure/storage-blob"
import Resizer from "react-image-file-resizer";

var signInOptions = {
    clientId: msalConfig.auth.clientId,
    tenantId: apiConfig.tenantId
}

var browserCredential = new InteractiveBrowserCredential(signInOptions);

const resizeFile = (file, maxWidth, maxHeight, format, quality, rotation) =>
    new Promise((resolve) => {
        Resizer.imageFileResizer(
            file,
            maxWidth,
            maxHeight,
            format,
            quality,
            rotation,
            (uri) => {
                resolve(uri);
            },
            "blob"
        );
    });

export async function uploadAndSetTags(containerName, files, collection, album) {
    if (files.length > 0) {
        files.forEach(async (file) => {
            var imageUrl = `${apiConfig.storageApiEndpoint}/${containerName}/${collection}/${album}/${file.file.name}`
            var thumbUrl = `${apiConfig.storageApiEndpoint}/${containerName}/${collection}/${album}/thumbs/${file.file.name}`

            var metadata = {
                exifData: file.exif
            }

            var tags = {
                collection: collection,
                album: album,
                name: file.file.name,
                size: String(file.file.size),
                url: imageUrl,
                thumbUrl: thumbUrl,
                isThumb: 'false'
            }
            
            // upload main image file
            let thumbBlockClient = new BlockBlobClient(imageUrl, browserCredential);
            await thumbBlockClient.upload(file.file, file.file.size);
            await thumbBlockClient.setTags(tags);
            await thumbBlockClient.setMetadata(metadata);

            // create & upload thumbnail
            tags.isThumb = 'true'
            const thumb = await resizeFile(file.file, 300, 300, "JPEG", 50, 0);
            let thumbFile = new File([thumb], file.file.name)
            let imageBlockClient = new BlockBlobClient(thumbUrl, browserCredential);
            await imageBlockClient.upload(thumbFile, thumbFile.size);
            await imageBlockClient.setTags(tags);

        })
    } else {
        console.log("no files found...")
    }
}