import { apiConfig, msalConfig } from "../authConfig";
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
            // replace unsupported characters
            let fileName = file.file.name.replace(/[^0-9a-zA-Z+-.:=_/]/g, '_');
            let image = new Image();
            image.src = file.file;

            var imageUrl = `${apiConfig.storageApiEndpoint}/${containerName}/${collection}/${album}/${fileName}`
            var thumbUrl = `${apiConfig.storageApiEndpoint}/thumbs/${collection}/${album}/${fileName}`

            if (file.exif.MakerNote) {
                file.exif.MakerNote = [];
            }

            var metadata = {
                exifData: JSON.stringify(file.exif),
                collection: collection,
                album: album,
                name: fileName,
                originalFileName: file.file.name
            }

            var tags = {
                collection: collection,
                album: album,
                name: fileName,
                size: String(file.file.size),
                url: imageUrl,
                thumbUrl: thumbUrl,
                isThumb: 'false',
                isAlbumThumb: file.isAlbumThumb
            }

            console.log("tags: " + JSON.stringify(tags));
            console.log(`metadata: ${JSON.stringify(metadata)}`)

            // upload main image file
            let thumbBlockClient = new BlockBlobClient(imageUrl, browserCredential);
            await thumbBlockClient.upload(file.file, file.file.size);
            await thumbBlockClient.setTags(tags);
            await thumbBlockClient.setMetadata(metadata);

            // create & upload thumbnail
            tags.isThumb = 'true'
            const thumb = await resizeFile(file.file, 300, 300, "JPEG", 50, 0);
            let thumbFile = new File([thumb], file.file.name)

            image = new Image();
            image.src = thumb;

            //tags.width = image.width;
            //tags.height = image.height;

            let imageBlockClient = new BlockBlobClient(thumbUrl, browserCredential);
            await imageBlockClient.upload(thumbFile, thumbFile.size);
            await imageBlockClient.setTags(tags);
            await imageBlockClient.setMetadata(metadata);
        })
    } else {
        console.log("no files found...")
    }
}