import { apiConfig, msalConfig } from "./authConfig";
import { InteractiveBrowserCredential } from "@azure/identity"
import { BlobServiceClient } from "@azure/storage-blob"

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

export async function getBlobsByTags(containerName, fileName, album, collection) {
    let tags = {
        collection: collection,
        album: album
    }

    let imageUrl = apiConfig.storageApiEndpoint + "/" + containerName + "/" + collection + "/" + album + "/" + fileName
    let blobServiceClient = new BlobServiceClient(imageUrl, browserCredential)
    let blobs = blobServiceClient.findBlobsByTags(tags);
    return blobs
}

export async function getAlbumCollections(containerName) {
    let blobServiceClient = new BlobServiceClient(apiConfig.storageApiEndpoint, browserCredential)
    let containerClient = blobServiceClient.getContainerClient(containerName);
    let albumCollectionMap = new Map();

    for await (const blob of containerClient.listBlobsFlat(options)) {
        if (blob.tags["collection"] !== undefined && (blob.tags["album"] !== undefined)) {
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


export async function getCollections(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.collectionApiEndpoint, options)
        .then(response => response.json())
        .catch(error => console.log(error));
}

export async function getAlbums(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.albumApiEndpoint, options)
        .then(response => response.json())
        .catch(error => console.log(error));
}

export async function getAlbumPhotos(accessToken) {
    const headers = new Headers();
    const bearer = `Bearer ${accessToken}`;

    headers.append("Authorization", bearer);
    headers.append("Access-Control-Allow-Headers", "*");
    headers.append("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE");
    headers.append("Access-Control-Allow-Origin", "http://localhost:3000");

    const options = {
        method: "GET",
        headers: headers,
    };

    return fetch(apiConfig.photoApiEndpoint + "?collection=2022&album=sport", options)
        .then(response => response.json())
        .catch(error => console.log(error));
}
