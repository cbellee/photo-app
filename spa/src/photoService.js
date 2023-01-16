import { apiConfig, msalConfig } from "./authConfig";
import { InteractiveBrowserCredential } from "@azure/identity"
import { BlobClient, BlobServiceClient, BlockBlobClient } from "@azure/storage-blob"

var signInOptions = {
    clientId: msalConfig.auth.clientId,
    tenantId: apiConfig.tenantId
}

var browserCredential = new InteractiveBrowserCredential(signInOptions);

const options = {
    includeDeleted: false,
    includeMetadata: true,
    includeTags: true
}

export async function getBlobsByTags(containerName, collection, album, thumbsOnly = true) {
    let tagQuery = `@container='${containerName}' AND collection='${collection}' AND album='${album}' AND isThumb='${thumbsOnly}'`;
    let i = 1;
    let blobs = [];

    const listOptions = {
        includeMetadata: true,
        includeSnapshots: false,
        includeTags: true,
        includeVersions: false
    };

    let storageUrl = apiConfig.storageApiEndpoint
    let blobServiceClient = new BlobServiceClient(storageUrl, browserCredential)
    for await (const blob of blobServiceClient.findBlobsByTags(tagQuery, listOptions)) {
        let containerClient = blobServiceClient.getContainerClient(blob.containerName);
        let blockBlobClient = containerClient.getBlobClient(blob.name);
        let tags = (await blockBlobClient.getTags()).tags;
        let metadata = (await blockBlobClient.getProperties());
        let b = {
            name: blockBlobClient.name,
            size: metadata.contentLength,
            createdOn: metadata.createdOn,
            lastModified: metadata.lastModified,
            metadata: metadata.metadata,
            tags: tags
        }
        blobs.push(b);
    }
    console.log("blob: " + JSON.stringify(blobs));
    return blobs;
}

export async function getAlbumCollections(containerName) {
    let blobServiceClient = new BlobServiceClient(apiConfig.storageApiEndpoint, browserCredential)
    let containerClient = blobServiceClient.getContainerClient(containerName);
    let albumCollectionMap = new Map();

    for await (const blob of containerClient.listBlobsFlat(options)) {
        if (blob.tags && blob.tags["collection"] !== undefined && (blob.tags["album"] !== undefined)) {
            let collection = blob.tags["collection"];
            let album = blob.tags["album"];

            if (!albumCollectionMap.has(collection)) {
                albumCollectionMap.set(collection, [album]);
            }

            if (!albumCollectionMap.get(collection).includes(album)) {
                albumCollectionMap.get(collection).push(album);
            }
        }
    }
    return albumCollectionMap;
}